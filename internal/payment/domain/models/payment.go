package models

import "eshop-monolith/pkg/utils"

// Payment 支付记录
type Payment struct {
	ID             int64   `json:"id"`
	OrderID        int64   `json:"order_id"`
	OrderType      string  `json:"order_type"` // 订单类型: "order"(常规), "flash"(闪购)
	Amount         int64   `json:"amount"` // 金额（分）
	Currency       string  `json:"currency"`
	PaymentMethod  string  `json:"payment_method"` // 支付方式：alipay, wechat, bank, cash, flash
	TransactionID  string  `json:"transaction_id"` // 第三方交易ID
	Status         string  `json:"status"` // pending, processing, success, failed, refunded
	FailureReason  string  `json:"failure_reason"`
	Metadata       string  `json:"metadata"` // 额外元数据
	PaidAt         *utils.Timestamp `json:"paid_at"`

	CreatedAt      utils.Timestamp `json:"created_at"`
	UpdatedAt      utils.Timestamp `json:"updated_at"`

	// 关联
	Order          *Order `json:"order,omitempty"`
	Transactions   []PaymentTransaction `json:"transactions,omitempty"`
	Refunds        []Refund `json:"refunds,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

// PaymentMethod 支付方式
type PaymentMethod struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"` // 支付方式编码：alipay, wechat, bank, cash
	Name        string `json:"name"` // 支付方式名称
	Description string `json:"description"`
	Config      string `json:"config"` // 支付方式配置
	Status      int    `json:"status"` // 状态：1-启用，0-禁用

	CreatedAt   utils.Timestamp `json:"created_at"`
	UpdatedAt   utils.Timestamp `json:"updated_at"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}

// PaymentTransaction 支付交易记录
type PaymentTransaction struct {
	ID            int64  `json:"id"`
	PaymentID     int64  `json:"payment_id"`
	TransactionID string `json:"transaction_id"` // 第三方交易ID
	Amount        int64  `json:"amount"` // 交易金额（分）
	Type          string `json:"type"` // payment, refund, capture, void
	Status        string `json:"status"` // pending, success, failed
	ResponseData  string `json:"response_data"` // 第三方响应数据
	ErrorData     string `json:"error_data"` // 错误信息

	CreatedAt     utils.Timestamp `json:"created_at"`
	UpdatedAt     utils.Timestamp `json:"updated_at"`

	// 关联
	Payment       *Payment `json:"payment,omitempty"`
}

func (PaymentTransaction) TableName() string {
	return "payment_transactions"
}

// Refund 退款记录
type Refund struct {
	ID             int64  `json:"id"`
	PaymentID      int64  `json:"payment_id"`
	OrderID        int64  `json:"order_id"`
	RefundAmount   int64  `json:"refund_amount"` // 退款金额（分）
	RefundReason   string `json:"refund_reason"`
	TransactionID  string `json:"transaction_id"` // 第三方退款交易ID
	Status         string `json:"status"` // pending, processing, success, failed
	FailureReason  string `json:"failure_reason"`

	CreatedAt      utils.Timestamp `json:"created_at"`
	UpdatedAt      utils.Timestamp `json:"updated_at"`

	// 关联
	Payment        *Payment `json:"payment,omitempty"`
	Order          *Order `json:"order,omitempty"`
}

func (Refund) TableName() string {
	return "refunds"
}

// Order 订单关联（为了避免循环依赖，这里只定义需要的字段）
type Order struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func (Order) TableName() string {
	return "orders"
}
