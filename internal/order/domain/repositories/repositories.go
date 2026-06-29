package repositories

import (
	"context"
	"eshop-monolith/internal/order/domain/models"
)

type IorderRepository interface {
	FindByID(ctx context.Context, id int64) (*models.Order, error)
}

func NewOrderRepository(db interface{}) IorderRepository {
	return &orderRepository{}
}

type orderRepository struct{}

func (r *orderRepository) FindByID(ctx context.Context, id int64) (*models.Order, error) {
	return &models.Order{}, nil
}
