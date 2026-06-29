package user

type CreateAddressReq struct {
	Consignee string `json:"consignee" binding:"required,max=64"`
	Phone     string `json:"phone" binding:"required,max=20"`
	Country   string `json:"country" binding:"max=32"`
	Province  string `json:"province" binding:"required,max=32"`
	City      string `json:"city" binding:"required,max=32"`
	District  string `json:"district" binding:"required,max=32"`
	Detail    string `json:"detail" binding:"required,max=256"`
	ZipCode   string `json:"zip_code" binding:"max=10"`
	Tag       string `json:"tag" binding:"omitempty,oneof=home office company other"`
	IsDefault bool   `json:"is_default"`
}

type UpdateAddressReq struct {
	Consignee *string `json:"consignee" binding:"omitempty,max=64"`
	Phone     *string `json:"phone" binding:"omitempty,max=20"`
	Country   *string `json:"country" binding:"omitempty,max=32"`
	Province  *string `json:"province" binding:"omitempty,max=32"`
	City      *string `json:"city" binding:"omitempty,max=32"`
	District  *string `json:"district" binding:"omitempty,max=32"`
	Detail    *string `json:"detail" binding:"omitempty,max=256"`
	ZipCode   *string `json:"zip_code" binding:"omitempty,max=10"`
	Tag       *string `json:"tag" binding:"omitempty,oneof=home office company other"`
	IsDefault *bool   `json:"is_default"`
}
