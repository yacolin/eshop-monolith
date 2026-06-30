package trade

import (
	"time"

	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type Order struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo       string         `gorm:"type:varchar(32);not null;uniqueIndex:uk_order_no" json:"order_no"`
	UserID        int64          `gorm:"not null;index:idx_user_id" json:"user_id"`
	TotalAmount   int64          `gorm:"not null;default:0" json:"total_amount"`
	DiscountAmount int64         `gorm:"not null;default:0" json:"discount_amount"`
	ShippingFee   int64          `gorm:"not null;default:0" json:"shipping_fee"`
	PayAmount     int64          `gorm:"not null;default:0" json:"pay_amount"`
	Status        string         `gorm:"type:varchar(20);not null;default:'pending';index:idx_status" json:"status"`
	PaymentStatus string         `gorm:"type:varchar(20);not null;default:'unpaid';index:idx_payment_status" json:"payment_status"`
	PaymentMethod string         `gorm:"type:varchar(32);default:''" json:"payment_method"`
	Consignee     string         `gorm:"type:varchar(64);not null;default:''" json:"consignee"`
	Phone         string         `gorm:"type:varchar(20);not null;default:''" json:"phone"`
	Province      string         `gorm:"type:varchar(32);default:''" json:"province"`
	City          string         `gorm:"type:varchar(32);default:''" json:"city"`
	District      string         `gorm:"type:varchar(32);default:''" json:"district"`
	DetailAddr    string         `gorm:"type:varchar(256);default:''" json:"detail_addr"`
	ZipCode       string         `gorm:"type:varchar(10);default:''" json:"zip_code"`
	CouponID      *int64         `gorm:"index" json:"coupon_id"`
	CouponSnapshot string        `gorm:"type:json" json:"coupon_snapshot,omitempty"`
	BuyerRemark   string         `gorm:"type:varchar(500);default:''" json:"buyer_remark"`
	SellerRemark  string         `gorm:"type:varchar(500);default:''" json:"seller_remark"`
	Source        string         `gorm:"type:varchar(20);not null;default:'pc'" json:"source"`
	PaidAt        *time.Time     `gorm:"type:datetime" json:"paid_at"`
	ShippedAt     *time.Time     `gorm:"type:datetime" json:"shipped_at"`
	DeliveredAt   *time.Time     `gorm:"type:datetime" json:"delivered_at"`
	CompletedAt   *time.Time     `gorm:"type:datetime" json:"completed_at"`
	ClosedAt      *time.Time     `gorm:"type:datetime" json:"closed_at"`
	CreatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt utils.Timestamp      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (Order) TableName() string { return "tx_orders" }
