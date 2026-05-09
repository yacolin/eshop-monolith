package repositories

import (
	"context"
	userModels "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

type IuserRepository interface {
	Create(ctx context.Context, user *userModels.User) error
	GetByID(ctx context.Context, userID int64) (*userModels.User, error)
	GetByIDWithInfo(ctx context.Context, userID int64) (*userModels.User, error)
	Update(ctx context.Context, user *userModels.User) error
	Delete(ctx context.Context, userID int64) error
	List(ctx context.Context, limit, offset int) ([]userModels.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IuserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *userModels.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		po := models.UserFromDomain(user)
		if err := tx.Create(po).Error; err != nil {
			return err
		}
		user.ID = po.ID
		if user.UserInfo != nil {
			user.UserInfo.UserID = user.ID
			infoPO := models.UserInfoFromDomain(user.UserInfo)
			if err := tx.Create(infoPO).Error; err != nil {
				return err
			}
			user.UserInfo.ID = infoPO.ID
		}
		return nil
	})
}

func (r *userRepository) GetByID(ctx context.Context, userID int64) (*userModels.User, error) {
	var po models.UserPO
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *userRepository) GetByIDWithInfo(ctx context.Context, userID int64) (*userModels.User, error) {
	var po models.UserPO
	err := r.db.WithContext(ctx).Preload("UserInfo").Where("id = ?", userID).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *userRepository) Update(ctx context.Context, user *userModels.User) error {
	po := models.UserFromDomain(user)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *userRepository) Delete(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserInfoPO{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.UserPO{}, "id = ?", userID).Error
	})
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]userModels.User, int64, error) {
	var pos []models.UserPO
	query := r.db.WithContext(ctx).Model(&models.UserPO{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	users := make([]userModels.User, len(pos))
	for i, po := range pos {
		users[i] = *po.ToDomain()
	}
	return users, total, nil
}
