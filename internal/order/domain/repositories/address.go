package repositories

import (
	"context"

	addressModels "eshop-monolith/internal/address/domain/models"
	addressSvc "eshop-monolith/internal/address/service"
)

// IaddressForOrder 订单模块查询地址的防火墙接口
type IaddressForOrder interface {
	GetAddress(ctx context.Context, userID, addressID int64) (*addressModels.Address, error)
}

// AddressForOrderAdapter 适配器，调用 AddressService 获取地址
type AddressForOrderAdapter struct {
	addressSvc *addressSvc.AddressService
}

func NewAddressForOrderAdapter(addressSvc *addressSvc.AddressService) IaddressForOrder {
	return &AddressForOrderAdapter{addressSvc: addressSvc}
}

func (a *AddressForOrderAdapter) GetAddress(ctx context.Context, userID, addressID int64) (*addressModels.Address, error) {
	return a.addressSvc.GetByID(ctx, userID, addressID)
}
