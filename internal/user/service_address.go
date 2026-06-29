package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
)

const maxAddressPerUser = 20

type AddressService struct {
	repo IaddressRepository
}

func NewAddressService(repo IaddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

func (s *AddressService) Create(ctx context.Context, userID int64, req *CreateAddressReq) (*Address, error) {
	count, err := s.repo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= maxAddressPerUser {
		return nil, errcode.ErrAddressLimit
	}

	addr := &Address{
		UserID:    userID,
		Consignee: req.Consignee,
		Phone:     req.Phone,
		Country:   req.Country,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		ZipCode:   req.ZipCode,
		Tag:       req.Tag,
		IsDefault: req.IsDefault,
	}

	if req.IsDefault {
		if err := s.repo.ClearDefaultByUserID(ctx, userID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Create(ctx, addr); err != nil {
		return nil, err
	}
	return addr, nil
}

type AddressListResult struct {
	Total int64      `json:"total"`
	List  []*Address `json:"list"`
}

func (s *AddressService) ListByUser(ctx context.Context, userID int64) (*AddressListResult, error) {
	list, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	total := int64(len(list))
	items := make([]*Address, total)
	for i := range list {
		items[i] = &list[i]
	}
	return &AddressListResult{Total: total, List: items}, nil
}

func (s *AddressService) GetByID(ctx context.Context, userID, addressID int64) (*Address, error) {
	addr, err := s.repo.FindByID(ctx, addressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrAddressNotFound
		}
		return nil, err
	}
	if addr.UserID != userID {
		return nil, errcode.ErrAddressNotFound
	}
	return addr, nil
}

func (s *AddressService) Update(ctx context.Context, userID int64, addressID int64, req *UpdateAddressReq) (*Address, error) {
	addr, err := s.GetByID(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}

	if req.Consignee != nil {
		addr.Consignee = *req.Consignee
	}
	if req.Phone != nil {
		addr.Phone = *req.Phone
	}
	if req.Country != nil {
		addr.Country = *req.Country
	}
	if req.Province != nil {
		addr.Province = *req.Province
	}
	if req.City != nil {
		addr.City = *req.City
	}
	if req.District != nil {
		addr.District = *req.District
	}
	if req.Detail != nil {
		addr.Detail = *req.Detail
	}
	if req.ZipCode != nil {
		addr.ZipCode = *req.ZipCode
	}
	if req.Tag != nil {
		addr.Tag = *req.Tag
	}

	if req.IsDefault != nil && *req.IsDefault && !addr.IsDefault {
		if err := s.repo.ClearDefaultByUserID(ctx, userID); err != nil {
			return nil, err
		}
		addr.IsDefault = true
	} else if req.IsDefault != nil {
		addr.IsDefault = *req.IsDefault
	}

	if err := s.repo.Update(ctx, addr); err != nil {
		return nil, err
	}
	return addr, nil
}

func (s *AddressService) Delete(ctx context.Context, userID int64, addressID int64) error {
	addr, err := s.GetByID(ctx, userID, addressID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, addr.ID)
}

func (s *AddressService) GetDefault(ctx context.Context, userID int64) (*Address, error) {
	addr, err := s.repo.GetDefaultByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrAddressNotFound
		}
		return nil, err
	}
	return addr, nil
}
