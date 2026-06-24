package models

import "eshop-monolith/pkg/utils"

// Address 收货地址
type Address struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Consignee string          `json:"consignee"`
	Phone     string          `json:"phone"`
	Province  string          `json:"province"`
	City      string          `json:"city"`
	District  string          `json:"district"`
	Detail    string          `json:"detail"`
	ZipCode   string          `json:"zip_code"`
	IsDefault bool            `json:"is_default"`
	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

func (Address) TableName() string { return "addresses" }
