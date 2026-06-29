package trade

// CreatePaymentReq 创建支付请求
type CreatePaymentReq struct {
	OrderNo       string `json:"order_no" binding:"required,max=32"`
	OrderID       int64  `json:"order_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	PaymentMethod string `json:"payment_method" binding:"required,max=32"`
	Channel       string `json:"channel" binding:"max=32"`
	OrderType     string `json:"order_type" binding:"max=20"`
}

// PaymentCallbackReq 支付回调请求（渠道通知）
type PaymentCallbackReq struct {
	PaymentNo     string `json:"payment_no" binding:"required,max=32"`
	TransactionID string `json:"transaction_id" binding:"required,max=128"`
	Channel       string `json:"channel" binding:"max=32"`
	Status        string `json:"status" binding:"required,oneof=success failed"`
	FailureReason string `json:"failure_reason" binding:"max=500"`
	RawBody       string `json:"raw_body"`
}

// CreateRefundReq 创建退款请求
type CreateRefundReq struct {
	PaymentNo string `json:"payment_no" binding:"required,max=32"`
	Amount    int64  `json:"amount" binding:"required,gt=0"`
	Reason    string `json:"reason" binding:"max=500"`
}

// RefundCallbackReq 退款回调请求
type RefundCallbackReq struct {
	RefundNo          string `json:"refund_no" binding:"required,max=32"`
	ChannelTransactionID string `json:"channel_transaction_id" binding:"max=128"`
	Status            string `json:"status" binding:"required,oneof=success failed"`
	FailureReason     string `json:"failure_reason" binding:"max=500"`
}
