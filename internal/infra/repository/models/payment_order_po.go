package models

import (
	payDomain "eshop-monolith/internal/payment/domain/models"
)

// PaymentOrderPO 支付模块内部订单引用持久化对象
type PaymentOrderPO struct {
	ID     int64  `gorm:"primaryKey;autoIncrement"`
	Status string `gorm:"type:varchar(20);not null"`
}

func (PaymentOrderPO) TableName() string { return "orders" }

func (po *PaymentOrderPO) ToDomain() *payDomain.Order {
	return &payDomain.Order{
		ID:     po.ID,
		Status: po.Status,
	}
}

func PaymentOrderFromDomain(o *payDomain.Order) *PaymentOrderPO {
	return &PaymentOrderPO{
		ID:     o.ID,
		Status: o.Status,
	}
}
