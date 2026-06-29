package review

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

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

func (s *ReviewService) CreateReview(ctx context.Context, in CreateReviewInput) (*Review, error) {
	if in.OverallRating < 1 || in.OverallRating > 5 {
		return nil, errcode.ErrReviewInvalidRating
	}
	if err := s.verifyPurchase(ctx, in); err != nil {
		return nil, err
	}

	rv := &Review{
		UserID:          in.UserID,
		OrderID:         in.OrderID,
		OrderItemID:     in.OrderItemID,
		SpuID:           in.SpuID,
		SkuID:           in.SkuID,
		OverallRating:   in.OverallRating,
		QualityRating:   in.QualityRating,
		LogisticsRating: in.LogisticsRating,
		ServiceRating:   in.ServiceRating,
		Content:         in.Content,
		IsAnonymous:     in.IsAnonymous,
		Status:          ReviewStatusPending,
	}
	if err := s.repo.CreateReview(ctx, rv); err != nil {
		return nil, fmt.Errorf("create review failed: %w", err)
	}
	return rv, nil
}

func (s *ReviewService) verifyPurchase(ctx context.Context, in CreateReviewInput) error {
	if s.findOrderByItem == nil || in.OrderItemID == nil {
		return nil
	}
	order, err := s.findOrderByItem(ctx, *in.OrderItemID)
	if err != nil || order == nil {
		return errcode.ErrReviewNotPurchased
	}
	if order.CustomerID != strconv.FormatInt(in.UserID, 10) {
		return errcode.ErrReviewNotPurchased
	}
	matched := false
	for _, item := range order.Items {
		if item.ID != *in.OrderItemID {
			continue
		}
		pid, perr := strconv.ParseInt(item.ProductID, 10, 64)
		if perr == nil && pid == in.SpuID {
			matched = true
			break
		}
	}
	if !matched {
		return errcode.ErrReviewNotPurchased
	}
	return nil
}

func (s *ReviewService) ListBySpu(ctx context.Context, spuID int64, statuses []int8, page, size int) ([]*Review, int64, error) {
	page, size = normalizePaging(page, size)
	list, total, err := s.repo.ListBySpu(ctx, spuID, statuses, page, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result, total, nil
}

func (s *ReviewService) ListByUser(ctx context.Context, userID int64, page, size int) ([]*Review, int64, error) {
	page, size = normalizePaging(page, size)
	list, total, err := s.repo.ListByUser(ctx, userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result, total, nil
}

func (s *ReviewService) ListPending(ctx context.Context, page, size int) ([]*Review, int64, error) {
	page, size = normalizePaging(page, size)
	list, total, err := s.repo.ListByStatus(ctx, ReviewStatusPending, page, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Review, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result, total, nil
}

func (s *ReviewService) GetByID(ctx context.Context, reviewID int64) (*Review, error) {
	rv, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrReviewNotFound
		}
		return nil, err
	}
	return rv, nil
}

func (s *ReviewService) ModerateReview(ctx context.Context, reviewID int64, auditorID int64, status int8, reason string) error {
	rv, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrReviewNotFound
		}
		return err
	}
	if err := s.repo.UpdateStatus(ctx, reviewID, status, reason); err != nil {
		return err
	}
	_ = s.repo.CreateAuditLog(ctx, &ReviewAuditLog{
		ReviewID:  reviewID,
		AuditorID: auditorID,
		Action:    auditActionFromStatus(status),
		Reason:    reason,
	})
	s.recomputeRating(ctx, rv.SpuID)
	return nil
}

func auditActionFromStatus(status int8) int8 {
	switch status {
	case ReviewStatusApproved:
		return AuditActionApprove
	case ReviewStatusRejected:
		return AuditActionReject
	default:
		return AuditActionReview
	}
}

func (s *ReviewService) ReplyReview(ctx context.Context, reviewID int64, replyBy int64, content string) (*ReviewReply, error) {
	rv, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrReviewNotFound
		}
		return nil, err
	}

	reply := &ReviewReply{
		ReviewID:    reviewID,
		ReplyBy:     replyBy,
		ReplyByType: ReplyByMerchant,
		Content:     content,
	}
	reply, err = s.repo.CreateReply(ctx, reply)
	if err != nil {
		return nil, err
	}

	newCount := rv.ReplyCount + 1
	if err := s.repo.UpdateLatestReply(ctx, reviewID, reply.ID, newCount); err != nil {
		return nil, err
	}
	return reply, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, userID int64, reviewID int64, isAdmin bool) error {
	rv, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrReviewNotFound
		}
		return err
	}
	if !isAdmin && rv.UserID != userID {
		return errcode.ErrReviewNotOwner
	}
	if err := s.repo.Delete(ctx, reviewID); err != nil {
		return err
	}
	s.recomputeRating(ctx, rv.SpuID)
	return nil
}

func (s *ReviewService) GetRatingSummary(ctx context.Context, spuID int64) (*ReviewRating, error) {
	return s.repo.GetRatingSummary(ctx, spuID)
}

func (s *ReviewService) recomputeRating(ctx context.Context, spuID int64) {
	counts, total, err := s.repo.CountRatingBySpu(ctx, spuID)
	if err != nil {
		return
	}
	_ = s.repo.UpsertRatingSummary(ctx, ComputeRatingSummary(spuID, counts, total))
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
