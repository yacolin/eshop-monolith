package service

import (
	"context"
	"fmt"
	"strconv"

	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/logger"
	"eshop-monolith/pkg/query"

	"eshop-monolith/internal/review/domain/models"
	"eshop-monolith/internal/review/domain/repositories"
	reviewEvents "eshop-monolith/internal/review/events"
)

// 常量
const (
	maxMediaCount = 9 // 单条评论最多媒体数量
)

// ---------------------------------------------------------------------------
// 跨模块依赖的适配（端口适配器模式）
//
// 订单/用户信息仓储定义在各自模块的 domain/repositories 包中，
// 且其接口（如 IorderRepository）在该包外不可见（小写开头）。
// 这里用本地快照 + 查找闭包做鸭子类型，仅声明评论模块所需的方法，
// 避免反向依赖其他模块的领域模型，保持模块边界清晰、不产生循环依赖。
// ---------------------------------------------------------------------------

// OrderSnapshot 订单快照（与 order 模块解耦）
type OrderSnapshot struct {
	ID         int64
	CustomerID string
	OrderNo    string
	Items      []OrderItemSnapshot
}

// OrderItemSnapshot 订单项快照
type OrderItemSnapshot struct {
	ID        int64
	ProductID string
}

// UserInfoSnapshot 用户信息快照（与 user 模块解耦）
type UserInfoSnapshot struct {
	Nickname string
	Avatar   string
}

// OrderByItemLookup 根据「订单项 ID」查询所属订单快照（在 routes 层适配 order 模块）
type OrderByItemLookup func(ctx context.Context, orderItemID int64) (*OrderSnapshot, error)

// UserInfoLookup 根据「用户 ID」查询用户信息快照（在 routes 层适配 user 模块）
type UserInfoLookup func(ctx context.Context, userID int64) (*UserInfoSnapshot, error)

// ReviewService 评论服务
type ReviewService struct {
	repo        *repositories.ReviewRepository
	bus         *eventbus.Bus
	findOrderByItem OrderByItemLookup
	findUser    UserInfoLookup
}

// NewReviewService 创建评论服务
//
// findOrderByItem: 将 order 模块的 *orderModels.Order 转换为本模块快照（按订单项反查订单）
// findUser: 将 user 模块的 *userModels.UserInfo 转换为本模块快照
func NewReviewService(
	repo *repositories.ReviewRepository,
	bus *eventbus.Bus,
	findOrderByItem OrderByItemLookup,
	findUser UserInfoLookup,
) *ReviewService {
	return &ReviewService{
		repo:            repo,
		bus:             bus,
		findOrderByItem: findOrderByItem,
		findUser:        findUser,
	}
}

// CreateReviewInput 创建评论入参
type CreateReviewInput struct {
	ProductID   int64
	UserID      int64
	OrderItemID int64
	Rating      int
	Content     string
	Media       []models.ReviewMedia
}

// CreateReview 创建评论
//
// 校验链：
//  1. 评分范围合法（1-5）
//  2. 媒体数量不超限
//  3. 该订单项尚未被评论过（防重复）
//  4. 该订单项所属订单属于当前用户（防越权）
//  5. 该订单项确实属于目标产品（防错评）
//
// 新评论默认进入「待审核」状态（pending），由管理端审核通过后对外可见。
func (s *ReviewService) CreateReview(ctx context.Context, in CreateReviewInput) (*models.Review, error) {
	// 1. 评分范围
	if in.Rating < 1 || in.Rating > 5 {
		return nil, errcode.ErrReviewInvalidRating
	}
	// 2. 媒体数量
	if len(in.Media) > maxMediaCount {
		return nil, errcode.ErrReviewMediaLimitExceed
	}
	// 3. 防重复评论
	exists, err := s.repo.ExistsByOrderItem(ctx, in.OrderItemID)
	if err != nil {
		return nil, fmt.Errorf("校验重复评论失败: %w", err)
	}
	if exists {
		return nil, errcode.ErrReviewDuplicate
	}
	// 4 & 5. 校验购买关系
	if err := s.verifyPurchase(ctx, in); err != nil {
		return nil, err
	}

	rv := &models.Review{
		ProductID:   in.ProductID,
		UserID:      in.UserID,
		OrderItemID: in.OrderItemID,
		Rating:      in.Rating,
		Content:     in.Content,
		Media:       in.Media,
		Status:      models.ReviewStatusPending,
	}
	if err := s.repo.Create(ctx, rv); err != nil {
		return nil, fmt.Errorf("创建评论失败: %w", err)
	}

	// 发布评论创建事件（异步，不阻塞主流程；推荐/搜索系统可订阅）
	if s.bus != nil {
		s.bus.PublishAsync(reviewEvents.ReviewCreatedEvent{
			ReviewID:  rv.ID,
			ProductID: rv.ProductID,
			UserID:    rv.UserID,
			Rating:    rv.Rating,
			Status:    string(rv.Status),
		})
	}
	return rv, nil
}

// verifyPurchase 校验用户购买关系
func (s *ReviewService) verifyPurchase(ctx context.Context, in CreateReviewInput) error {
	if s.findOrderByItem == nil {
		// 未配置订单适配器时降级：记录告警并放行（仅在非生产/测试场景）
		logger.Warn("order lookup not configured, skip purchase verification",
			"product_id", in.ProductID, "order_item_id", in.OrderItemID)
		return nil
	}

	order, err := s.findOrderByItem(ctx, in.OrderItemID)
	if err != nil || order == nil {
		return errcode.ErrReviewNotPurchased
	}

	// 订单归属校验：order.CustomerID 为字符串型用户 ID
	if order.CustomerID != strconv.FormatInt(in.UserID, 10) {
		return errcode.ErrReviewNotPurchased
	}

	// 订单项与产品匹配校验
	matched := false
	for _, item := range order.Items {
		if item.ID != in.OrderItemID {
			continue
		}
		pid, perr := strconv.ParseInt(item.ProductID, 10, 64)
		if perr == nil && pid == in.ProductID {
			matched = true
			break
		}
	}
	if !matched {
		return errcode.ErrReviewNotPurchased
	}
	return nil
}

// ListProductReviews 查询产品可见评论（默认仅 approved），并附带评论者昵称/头像
func (s *ReviewService) ListProductReviews(ctx context.Context, productID int64, statuses []models.ReviewStatus, page, size int) (*query.ListResult[*models.Review], error) {
	page, size = normalizePaging(page, size)

	strStatuses := make([]string, 0, len(statuses))
	for _, st := range statuses {
		strStatuses = append(strStatuses, string(st))
	}

	list, total, err := s.repo.ListByProduct(ctx, productID, strStatuses, page, size)
	if err != nil {
		return nil, fmt.Errorf("查询产品评论失败: %w", err)
	}

	result := make([]*models.Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	s.enrichReviewers(ctx, result)

	return &query.ListResult[*models.Review]{Total: total, List: result}, nil
}

// ListUserReviews 查询用户自己的评论（含所有状态）
func (s *ReviewService) ListUserReviews(ctx context.Context, userID int64, page, size int) (*query.ListResult[*models.Review], error) {
	page, size = normalizePaging(page, size)

	list, total, err := s.repo.ListByUser(ctx, userID, page, size)
	if err != nil {
		return nil, fmt.Errorf("查询用户评论失败: %w", err)
	}

	result := make([]*models.Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	s.enrichReviewers(ctx, result)

	return &query.ListResult[*models.Review]{Total: total, List: result}, nil
}

// ListPendingReviews 查询待审核评论（管理端）
func (s *ReviewService) ListPendingReviews(ctx context.Context, page, size int) (*query.ListResult[*models.Review], error) {
	page, size = normalizePaging(page, size)

	list, total, err := s.repo.ListByStatus(ctx, string(models.ReviewStatusPending), page, size)
	if err != nil {
		return nil, fmt.Errorf("查询待审核评论失败: %w", err)
	}

	result := make([]*models.Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	s.enrichReviewers(ctx, result)

	return &query.ListResult[*models.Review]{Total: total, List: result}, nil
}

// ModerateReview 审核评论（通过/拒绝/隐藏），并刷新产品评分汇总
func (s *ReviewService) ModerateReview(ctx context.Context, reviewID int64, status models.ReviewStatus) error {
	rv, err := s.repo.GetByID(ctx, reviewID)
	if err != nil {
		return errcode.ErrReviewNotFound
	}
	if err := s.repo.UpdateStatus(ctx, reviewID, status); err != nil {
		return fmt.Errorf("更新评论状态失败: %w", err)
	}

	// 审核状态变化会改变对外可见的评论数，需重算评分汇总
	s.recomputeRatingSummary(ctx, rv.ProductID)

	if s.bus != nil {
		s.bus.PublishAsync(reviewEvents.ReviewModeratedEvent{
			ReviewID:  reviewID,
			ProductID: rv.ProductID,
			Status:    string(status),
		})
	}
	return nil
}

// ReplyReview 商家回复评论
func (s *ReviewService) ReplyReview(ctx context.Context, reviewID int64, reply string) error {
	if _, err := s.repo.GetByID(ctx, reviewID); err != nil {
		return errcode.ErrReviewNotFound
	}
	return s.repo.UpdateReply(ctx, reviewID, reply)
}

// GetReview 获取单条评论（用于归属校验等）
func (s *ReviewService) GetReview(ctx context.Context, reviewID int64) (*models.Review, error) {
	rv, err := s.repo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, errcode.ErrReviewNotFound
	}
	return rv, nil
}

// DeleteReview 删除评论（评论者本人或管理员），并刷新评分汇总
func (s *ReviewService) DeleteReview(ctx context.Context, reviewID int64) error {
	rv, err := s.repo.GetByID(ctx, reviewID)
	if err != nil {
		return errcode.ErrReviewNotFound
	}
	if err := s.repo.Delete(ctx, reviewID); err != nil {
		return fmt.Errorf("删除评论失败: %w", err)
	}

	s.recomputeRatingSummary(ctx, rv.ProductID)

	if s.bus != nil {
		s.bus.PublishAsync(reviewEvents.ReviewDeletedEvent{
			ReviewID:  reviewID,
			ProductID: rv.ProductID,
		})
	}
	return nil
}

// GetProductRating 获取产品评分汇总
func (s *ReviewService) GetProductRating(ctx context.Context, productID int64) (*models.ProductRatingSummary, error) {
	return s.repo.GetRatingSummary(ctx, productID)
}

// recomputeRatingSummary 重算并持久化产品评分汇总（仅统计 approved 评论）
func (s *ReviewService) recomputeRatingSummary(ctx context.Context, productID int64) {
	stats, err := s.repo.CountRatingByProduct(ctx, productID)
	if err != nil {
		logger.Error("统计产品评分失败", "product_id", productID, "error", err)
		return
	}
	summary, err := s.repo.UpsertRatingSummary(ctx, productID, stats)
	if err != nil {
		logger.Error("更新评分汇总失败", "product_id", productID, "error", err)
		return
	}
	if s.bus != nil {
		s.bus.PublishAsync(reviewEvents.RatingSummaryUpdatedEvent{
			ProductID:     summary.ProductID,
			AverageRating: summary.AverageRating,
			ReviewCount:   summary.ReviewCount,
		})
	}
}

// enrichReviewers 批量填充评论者昵称/头像
func (s *ReviewService) enrichReviewers(ctx context.Context, reviews []*models.Review) {
	if s.findUser == nil || len(reviews) == 0 {
		return
	}
	cache := make(map[int64]*UserInfoSnapshot, len(reviews))
	for _, rv := range reviews {
		if rv == nil {
			continue
		}
		if _, ok := cache[rv.UserID]; ok {
			continue
		}
		info, err := s.findUser(ctx, rv.UserID)
		if err != nil || info == nil {
			cache[rv.UserID] = nil
			continue
		}
		cache[rv.UserID] = info
	}
	for _, rv := range reviews {
		if rv == nil {
			continue
		}
		if info, ok := cache[rv.UserID]; ok && info != nil {
			rv.UserName = info.Nickname
			rv.UserAvatar = info.Avatar
		}
	}
}

// normalizePaging 规范化分页参数
func normalizePaging(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
