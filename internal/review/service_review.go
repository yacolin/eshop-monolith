package review

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"eshop-monolith/pkg/errcode"
)

const maxMediaCount = 9

type OrderSnapshot struct {
	ID         int64
	CustomerID string
	OrderNo    string
	Items      []OrderItemSnapshot
}

type OrderItemSnapshot struct {
	ID        int64
	ProductID string
}

type UserInfoSnapshot struct {
	Nickname string
	Avatar   string
}

type OrderByItemLookup func(ctx context.Context, orderItemID int64) (*OrderSnapshot, error)

type UserInfoLookup func(ctx context.Context, userID int64) (*UserInfoSnapshot, error)

type ReviewService struct {
	repo            IreviewRepository
	findOrderByItem OrderByItemLookup
	findUser        UserInfoLookup
}

func NewReviewService(repo IreviewRepository, findOrderByItem OrderByItemLookup, findUser UserInfoLookup) *ReviewService {
	return &ReviewService{
		repo:            repo,
		findOrderByItem: findOrderByItem,
		findUser:        findUser,
	}
}

type CreateReviewInput struct {
	ProductID   int64
	UserID      int64
	OrderItemID int64
	Rating      int
	Content     string
	Media       []ReviewMedia
}

func (s *ReviewService) CreateReview(ctx context.Context, in CreateReviewInput) (*Review, error) {
	if in.Rating < 1 || in.Rating > 5 {
		return nil, errcode.ErrReviewInvalidRating
	}
	if len(in.Media) > maxMediaCount {
		return nil, errcode.ErrReviewMediaLimitExceed
	}
	exists, err := s.repo.ExistsByOrderItem(ctx, in.OrderItemID)
	if err != nil {
		return nil, fmt.Errorf("check duplicate failed: %w", err)
	}
	if exists {
		return nil, errcode.ErrReviewDuplicate
	}
	if err := s.verifyPurchase(ctx, in); err != nil {
		return nil, err
	}

	mediaJSON, _ := json.Marshal(in.Media)
	rv := &Review{
		ProductID:   in.ProductID,
		UserID:      in.UserID,
		OrderItemID: in.OrderItemID,
		Rating:      in.Rating,
		Content:     in.Content,
		MediaJSON:   string(mediaJSON),
		Status:      string(ReviewStatusPending),
	}
	if err := s.repo.Create(ctx, rv); err != nil {
		return nil, fmt.Errorf("create review failed: %w", err)
	}
	return rv, nil
}

func (s *ReviewService) verifyPurchase(ctx context.Context, in CreateReviewInput) error {
	if s.findOrderByItem == nil {
		return nil
	}
	order, err := s.findOrderByItem(ctx, in.OrderItemID)
	if err != nil || order == nil {
		return errcode.ErrReviewNotPurchased
	}
	if order.CustomerID != strconv.FormatInt(in.UserID, 10) {
		return errcode.ErrReviewNotPurchased
	}
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

func (s *ReviewService) ListProductReviews(ctx context.Context, productID int64, statuses []string, page, size int) ([]*Review, int64, error) {
	page, size = normalizePaging(page, size)
	list, total, err := s.repo.ListByProduct(ctx, productID, statuses, page, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	s.enrichReviewers(ctx, result)
	return result, total, nil
}

func (s *ReviewService) ListUserReviews(ctx context.Context, userID int64, page, size int) ([]*Review, int64, error) {
	page, size = normalizePaging(page, size)
	list, total, err := s.repo.ListByUser(ctx, userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	s.enrichReviewers(ctx, result)
	return result, total, nil
}

func (s *ReviewService) ListPendingReviews(ctx context.Context, page, size int) ([]*Review, int64, error) {
	page, size = normalizePaging(page, size)
	list, total, err := s.repo.ListByStatus(ctx, string(ReviewStatusPending), page, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	s.enrichReviewers(ctx, result)
	return result, total, nil
}

func (s *ReviewService) ModerateReview(ctx context.Context, reviewID int64, status string) error {
	rv, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		return errcode.ErrReviewNotFound
	}
	if err := s.repo.UpdateStatus(ctx, reviewID, status); err != nil {
		return err
	}
	s.recomputeRatingSummary(ctx, rv.ProductID)
	return nil
}

func (s *ReviewService) ReplyReview(ctx context.Context, reviewID int64, reply string) error {
	if _, err := s.repo.FindByID(ctx, reviewID); err != nil {
		return errcode.ErrReviewNotFound
	}
	return s.repo.UpdateReply(ctx, reviewID, reply)
}

func (s *ReviewService) GetReview(ctx context.Context, reviewID int64) (*Review, error) {
	rv, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		return nil, errcode.ErrReviewNotFound
	}
	return rv, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewID int64) error {
	rv, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		return errcode.ErrReviewNotFound
	}
	if err := s.repo.Delete(ctx, reviewID); err != nil {
		return err
	}
	s.recomputeRatingSummary(ctx, rv.ProductID)
	return nil
}

func (s *ReviewService) GetProductRating(ctx context.Context, productID int64) (*ProductRatingSummary, error) {
	return s.repo.GetRatingSummary(ctx, productID)
}

func (s *ReviewService) recomputeRatingSummary(ctx context.Context, productID int64) {
	stats, err := s.repo.CountRatingByProduct(ctx, productID)
	if err != nil {
		return
	}
	_, _ = s.repo.UpsertRatingSummary(ctx, productID, stats)
}

func (s *ReviewService) enrichReviewers(ctx context.Context, reviews []*Review) {
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
