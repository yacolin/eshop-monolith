package service

import (
	"context"
	"eshop-monolith/internal/address/api/dto"
	"eshop-monolith/internal/address/domain/models"
	"eshop-monolith/internal/address/domain/repositories"
	"eshop-monolith/internal/address/events"
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/pkg/errcode"
	"gorm.io/gorm"
)

const maxAddressPerUser = 20

type AddressService struct {
	repo repositories.IaddressRepository
	db   *gorm.DB
	rabbit  *rabbitmq.Client
}

func NewAddressService(repo repositories.IaddressRepository, db *gorm.DB, rabbit *rabbitmq.Client) *AddressService {
	return &AddressService{repo: repo, db: db, rabbit: rabbit}
}

func (s *AddressService) Create(ctx context.Context, userID int64, req *dto.CreateAddressReq) (*models.Address, error) {
	count, err := s.repo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= maxAddressPerUser {
		return nil, errcode.ErrInvalidParams
	}

	addr := &models.Address{
		UserID: userID, Consignee: req.Consignee, Phone: req.Phone,
		Province: req.Province, City: req.City, District: req.District,
		Detail: req.Detail, ZipCode: req.ZipCode, IsDefault: req.IsDefault,
	}

	if req.IsDefault {
		if err := s.clearDefault(ctx, userID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Create(ctx, addr); err != nil {
		return nil, err
	}

	s.rabbit.Publish(context.Background(), events.AddressCreatedEvent{AddressID: addr.ID, UserID: userID})
	return addr, nil
}

func (s *AddressService) ListByUser(ctx context.Context, userID int64) (*dto.AddressListResp, error) {
	list, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	total := int64(len(list))
	resp := &dto.AddressListResp{Total: total, List: make([]*dto.AddressResp, total)}
	for i, a := range list {
		resp.List[i] = dto.ToAddressResp(&a)
	}
	return resp, nil
}

func (s *AddressService) GetByID(ctx context.Context, userID, addressID int64) (*models.Address, error) {
	addr, err := s.repo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if addr.UserID != userID {
		return nil, errcode.ErrNotFound
	}
	return addr, nil
}

func (s *AddressService) Update(ctx context.Context, userID int64, addressID int64, req *dto.UpdateAddressReq) (*models.Address, error) {
	addr, err := s.GetByID(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}

	if req.Consignee != nil { addr.Consignee = *req.Consignee }
	if req.Phone != nil { addr.Phone = *req.Phone }
	if req.Province != nil { addr.Province = *req.Province }
	if req.City != nil { addr.City = *req.City }
	if req.District != nil { addr.District = *req.District }
	if req.Detail != nil { addr.Detail = *req.Detail }
	if req.ZipCode != nil { addr.ZipCode = *req.ZipCode }

	if req.IsDefault != nil && *req.IsDefault && !addr.IsDefault {
		if err := s.clearDefault(ctx, userID); err != nil {
			return nil, err
		}
		addr.IsDefault = true
	} else if req.IsDefault != nil {
		addr.IsDefault = *req.IsDefault
	}

	if err := s.repo.Update(ctx, addr); err != nil {
		return nil, err
	}

	s.rabbit.Publish(context.Background(), events.AddressUpdatedEvent{AddressID: addr.ID, UserID: userID})
	return addr, nil
}

func (s *AddressService) Delete(ctx context.Context, userID int64, addressID int64) error {
	addr, err := s.GetByID(ctx, userID, addressID)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, addr.ID); err != nil {
		return err
	}
	s.rabbit.Publish(context.Background(), events.AddressDeletedEvent{AddressID: addressID, UserID: userID})
	return nil
}

func (s *AddressService) GetDefault(ctx context.Context, userID int64) (*models.Address, error) {
	addr, err := s.repo.GetDefaultByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func (s *AddressService) clearDefault(ctx context.Context, userID int64) error {
	return s.repo.ClearDefaultByUserID(ctx, userID)
}
