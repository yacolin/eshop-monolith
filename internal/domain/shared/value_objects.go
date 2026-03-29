package shared

// Money 金额值对象
type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// Address 地址值对象
type Address struct {
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
	ZipCode  string `json:"zip_code"`
}

// PhoneNumber 电话号码值对象
type PhoneNumber struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}

// Email 邮箱值对象
type Email struct {
	Value string `json:"value"`
}

// Password 密码值对象
type Password struct {
	Hash string `json:"hash"`
}

// ID 通用ID值对象
type ID struct {
	Value string `json:"value"`
}

// Quantity 数量值对象
type Quantity struct {
	Value int `json:"value"`
}

// SKU 商品SKU值对象
type SKU struct {
	Value string `json:"value"`
}