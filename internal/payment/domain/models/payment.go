package models

import (
	"eshop-monolith/internal/pkg/utils"

	"gorm.io/gorm"
)

// Payment 支付记录
type Payment struct {
	ID             int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID        int64  `json:"order_id" gorm:"not null;index"`
	Amount         int64  `json:"amount" gorm:"not null"` // 金额（分）
	Currency       string `json:"currency" gorm:"type:varchar(10);not null;default:'CNY'"`
	PaymentMethod  string `json:"payment_method" gorm:"type:varchar(50);not null"` // 支付方式：alipay, wechat, bank, cash
	TransactionID  string `json:"transaction_id" gorm:"type:varchar(255);index"` // 第三方交易ID
	Status         string `json:"status" gorm:"type:varchar(20);not null;default:'pending'"` // pending, processing, success, failed, refunded
	FailureReason  string `json:"failure_reason" gorm:"type:text"`
	Metadata       string `json:"metadata" gorm:"type:json"` // 额外元数据
	PaidAt         *utils.Timestamp `json:"paid_at" gorm:"type:timestamp"`

	CreatedAt      utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt      utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Order          *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Transactions   []PaymentTransaction `gorm:"foreignKey:PaymentID" json:"transactions,omitempty"`
	Refunds        []Refund `gorm:"foreignKey:PaymentID" json:"refunds,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

// PaymentMethod 支付方式
type PaymentMethod struct {
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string `json:"code" gorm:"type:varchar(50);not null;uniqueIndex"` // 支付方式编码：alipay, wechat, bank, cash
	Name        string `json:"name" gorm:"type:varchar(100);not null"` // 支付方式名称
	Description string `json:"description" gorm:"type:text"`
	Config      string `json:"config" gorm:"type:json"` // 支付方式配置
	Status      int    `json:"status" gorm:"type:tinyint;default:1"` // 状态：1-启用，0-禁用

	CreatedAt   utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}

// PaymentTransaction 支付交易记录
type PaymentTransaction struct {
	ID            int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentID     int64  `json:"payment_id" gorm:"not null;index"`
	TransactionID string `json:"transaction_id" gorm:"type:varchar(255);not null;index"` // 第三方交易ID
	Amount        int64  `json:"amount" gorm:"not null"` // 交易金额（分）
	Type          string `json:"type" gorm:"type:varchar(20);not null"` // payment, refund, capture, void
	Status        string `json:"status" gorm:"type:varchar(20);not null"` // pending, success, failed
	ResponseData  string `json:"response_data" gorm:"type:json"` // 第三方响应数据
	ErrorData     string `json:"error_data" gorm:"type:json"` // 错误信息

	CreatedAt     utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt     utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Payment       *Payment `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
}

func (PaymentTransaction) TableName() string {
	return "payment_transactions"
}

// Refund 退款记录
type Refund struct {
	ID             int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentID      int64  `json:"payment_id" gorm:"not null;index"`
	OrderID        int64  `json:"order_id" gorm:"not null;index"`
	RefundAmount   int64  `json:"refund_amount" gorm:"not null"` // 退款金额（分）
	RefundReason   string `json:"refund_reason" gorm:"type:text"`
	TransactionID  string `json:"transaction_id" gorm:"type:varchar(255);index"` // 第三方退款交易ID
	Status         string `json:"status" gorm:"type:varchar(20);not null;default:'pending'"` // pending, processing, success, failed
	FailureReason  string `json:"failure_reason" gorm:"type:text"`

	CreatedAt      utils.Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt      utils.Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Payment        *Payment `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	Order          *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (Refund) TableName() string {
	return "refunds"
}

// Order 订单关联（为了避免循环依赖，这里只定义需要的字段）
type Order struct {
	ID     int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Status string `json:"status" gorm:"type:varchar(20);not null"`
}

func (Order) TableName() string {
	return "orders"
}
