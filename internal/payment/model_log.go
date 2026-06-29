package payment

import "time"

type PaymentLog struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PaymentID     int64     `gorm:"not null;index:idx_payment_id" json:"payment_id"`
	PaymentNo     string    `gorm:"type:varchar(32);not null;index:idx_payment_no" json:"payment_no"`
	Channel       string    `gorm:"type:varchar(32);default:''" json:"channel"`
	TransactionID string    `gorm:"type:varchar(128);default:'';index:idx_transaction_id" json:"transaction_id"`
	Action        string    `gorm:"type:varchar(30);not null;index:idx_action" json:"action"`
	RequestBody   string    `gorm:"type:text" json:"request_body"`
	ResponseBody  string    `gorm:"type:text" json:"response_body"`
	Status        string    `gorm:"type:varchar(20);default:''" json:"status"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"created_at"`
}

func (PaymentLog) TableName() string { return "tx_payment_logs" }
