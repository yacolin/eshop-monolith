package dto

import (
	"eshop-monolith/internal/pkg/query"
)

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID       int64  `json:"order_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	ReturnURL     string `json:"return_url"`
	CallbackURL   string `json:"callback_url"`
	Metadata      string `json:"metadata"`
}

type PaymentListResult = query.ListResult[PaymentResponse]

// PaymentResponse 支付响应
type PaymentResponse struct {
	ID            int64  `json:"id"`
	OrderID       int64  `json:"order_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	PaymentURL    string `json:"payment_url,omitempty"`
	QRCode        string `json:"qrcode,omitempty"`
	CreatedAt     string `json:"created_at"`
	PaidAt        string `json:"paid_at,omitempty"`
}

// PaymentListQuery 支付列表查询参数
type PaymentListQuery struct {
	query.Pagination
	OrderID       int64  `form:"order_id"`
	PaymentMethod string `form:"payment_method"`
	Status        string `form:"status"`
	StartDate     string `form:"start_date"`
	EndDate       string `form:"end_date"`
	SortBy        string `form:"sort_by,default=created_at"`
	Order         string `form:"order,default=desc"`
}

// PaymentListResponse 支付列表响应
type PaymentListResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
}

// UpdatePaymentStatusRequest 更新支付状态请求
type UpdatePaymentStatusRequest struct {
	Status        string `json:"status" binding:"required"`
	TransactionID string `json:"transaction_id"`
	FailureReason string `json:"failure_reason"`
}

// CreateRefundRequest 创建退款请求
type CreateRefundRequest struct {
	PaymentID    int64  `json:"payment_id" binding:"required"`
	OrderID      int64  `json:"order_id" binding:"required"`
	RefundAmount int64  `json:"refund_amount" binding:"required,gt=0"`
	RefundReason string `json:"refund_reason" binding:"required"`
}

// RefundResponse 退款响应
type RefundResponse struct {
	ID            int64  `json:"id"`
	PaymentID     int64  `json:"payment_id"`
	OrderID       int64  `json:"order_id"`
	RefundAmount  int64  `json:"refund_amount"`
	RefundReason  string `json:"refund_reason"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// RefundListQuery 退款列表查询参数
type RefundListQuery struct {
	PaymentID int64  `form:"payment_id"`
	OrderID   int64  `form:"order_id"`
	Status    string `form:"status"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
	SortBy    string `form:"sort_by,default=created_at"`
	Order     string `form:"order,default=desc"`
}

// RefundListResult 退款列表结果
type RefundListResult = query.ListResult[RefundResponse]

// PaymentMethodResponse 支付方式响应
type PaymentMethodResponse struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// PaymentMethodListResponse 支付方式列表响应
type PaymentMethodListResponse struct {
	PaymentMethods []PaymentMethodResponse `json:"payment_methods"`
	Total          int64                   `json:"total"`
	Page           int                     `json:"page"`
	PageSize       int                     `json:"page_size"`
}
