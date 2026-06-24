package dto

import (
	"eshop-monolith/internal/address/domain/models"
	"eshop-monolith/pkg/utils"
)

type CreateAddressReq struct {
	Consignee string `json:"consignee" binding:"required,max=64"`
	Phone     string `json:"phone" binding:"required,max=20"`
	Province  string `json:"province" binding:"required,max=32"`
	City      string `json:"city" binding:"required,max=32"`
	District  string `json:"district" binding:"required,max=32"`
	Detail    string `json:"detail" binding:"required,max=256"`
	ZipCode   string `json:"zip_code" binding:"max=10"`
	IsDefault bool   `json:"is_default"`
}

type UpdateAddressReq struct {
	Consignee *string `json:"consignee" binding:"omitempty,max=64"`
	Phone     *string `json:"phone" binding:"omitempty,max=20"`
	Province  *string `json:"province" binding:"omitempty,max=32"`
	City      *string `json:"city" binding:"omitempty,max=32"`
	District  *string `json:"district" binding:"omitempty,max=32"`
	Detail    *string `json:"detail" binding:"omitempty,max=256"`
	ZipCode   *string `json:"zip_code" binding:"omitempty,max=10"`
	IsDefault *bool   `json:"is_default"`
}

type AddressResp struct {
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

type AddressListResp struct {
	Total int64          `json:"total"`
	List  []*AddressResp `json:"list"`
}

func ToAddressResp(a *models.Address) *AddressResp {
	return &AddressResp{
		ID:        a.ID,
		UserID:    a.UserID,
		Consignee: a.Consignee,
		Phone:     a.Phone,
		Province:  a.Province,
		City:      a.City,
		District:  a.District,
		Detail:    a.Detail,
		ZipCode:   a.ZipCode,
		IsDefault: a.IsDefault,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
