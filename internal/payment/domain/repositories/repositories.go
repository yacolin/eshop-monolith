package repositories

import "gorm.io/gorm"

type IPaymentRepository interface {
	CreateWithTx(tx *gorm.DB, payment interface{}) error
}

func NewPaymentRepository(db interface{}) IPaymentRepository {
	return &paymentRepository{}
}

type paymentRepository struct{}

func (r *paymentRepository) CreateWithTx(tx *gorm.DB, payment interface{}) error {
	return nil
}
