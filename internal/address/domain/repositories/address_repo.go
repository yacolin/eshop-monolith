package repositories

import (
	"context"

	"eshop-monolith/internal/address/domain/models"
	repoModels "eshop-monolith/internal/infra/repository/models"
	"gorm.io/gorm"
)

// IaddressRepository 地址仓储接口
type IaddressRepository interface {
	Create(ctx context.Context, address *models.Address) error
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	ListByUserID(ctx context.Context, userID int64) ([]models.Address, error)
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id int64) error
	GetDefaultByUserID(ctx context.Context, userID int64) (*models.Address, error)
	CountByUserID(ctx context.Context, userID int64) (int64, error)
	ClearDefaultByUserID(ctx context.Context, userID int64) error
}

type addressRepository struct{ db *gorm.DB }

func NewAddressRepository(db *gorm.DB) IaddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *models.Address) error {
	po := repoModels.AddressFromDomain(address)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	address.ID = po.ID
	return nil
}

func (r *addressRepository) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	var po repoModels.AddressPO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *addressRepository) ListByUserID(ctx context.Context, userID int64) ([]models.Address, error) {
	var pos []repoModels.AddressPO
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&pos).Error; err != nil {
		return nil, err
	}
	addresses := make([]models.Address, len(pos))
	for i, po := range pos {
		addresses[i] = *po.ToDomain()
	}
	return addresses, nil
}

func (r *addressRepository) Update(ctx context.Context, address *models.Address) error {
	po := repoModels.AddressFromDomain(address)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *addressRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&repoModels.AddressPO{}, "id = ?", id).Error
}

func (r *addressRepository) GetDefaultByUserID(ctx context.Context, userID int64) (*models.Address, error) {
	var po repoModels.AddressPO
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_default = ?", userID, true).First(&po).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *addressRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&repoModels.AddressPO{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *addressRepository) ClearDefaultByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&repoModels.AddressPO{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error
}
