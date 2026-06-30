package trade

import "eshop-monolith/pkg/utils"

type OrderLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      int64     `gorm:"not null;index:idx_order_id" json:"order_id"`
	OrderNo      string    `gorm:"type:varchar(32);not null;index:idx_order_no" json:"order_no"`
	FromStatus   string    `gorm:"type:varchar(20);default:''" json:"from_status"`
	ToStatus     string    `gorm:"type:varchar(20);not null;index:idx_to_status" json:"to_status"`
	Operator     string    `gorm:"type:varchar(50);default:'system'" json:"operator"`
	OperatorType string    `gorm:"type:varchar(20);default:'system'" json:"operator_type"`
	Note         string    `gorm:"type:varchar(500);default:''" json:"note"`
	CreatedAt utils.Timestamp `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (OrderLog) TableName() string { return "tx_order_logs" }
