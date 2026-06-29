package user

import (
	"context"

	"gorm.io/gorm"
)

type IaddressRepository interface {
	Create(ctx context.Context, address *Address) error
	FindByID(ctx context.Context, id int64) (*Address, error)
	ListByUserID(ctx context.Context, userID int64) ([]Address, error)
	Update(ctx context.Context, address *Address) error
	Delete(ctx context.Context, id int64) error
	GetDefaultByUserID(ctx context.Context, userID int64) (*Address, error)
	CountByUserID(ctx context.Context, userID int64) (int64, error)
	ClearDefaultByUserID(ctx context.Context, userID int64) error
}

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) IaddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) Create(ctx context.Context, address *Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *AddressRepository) FindByID(ctx context.Context, id int64) (*Address, error) {
	var addr Address
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&addr).Error
	return &addr, err
}

func (r *AddressRepository) ListByUserID(ctx context.Context, userID int64) ([]Address, error) {
	var list []Address
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("is_default DESC, id DESC").Find(&list).Error
	return list, err
}

func (r *AddressRepository) Update(ctx context.Context, address *Address) error {
	return r.db.WithContext(ctx).Model(address).Updates(address).Error
}

func (r *AddressRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Address{}).Error
}

func (r *AddressRepository) GetDefaultByUserID(ctx context.Context, userID int64) (*Address, error) {
	var addr Address
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_default = ?", userID, true).
		First(&addr).Error
	return &addr, err
}

func (r *AddressRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Address{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *AddressRepository) ClearDefaultByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&Address{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error
}
