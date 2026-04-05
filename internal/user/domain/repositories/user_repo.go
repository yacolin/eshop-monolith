package repositories

import (
	"context"
	"eshop-monolith/internal/user/domain/models"

	"gorm.io/gorm"
)

type IuserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID int64) (*models.User, error)
	GetByIDWithInfo(ctx context.Context, userID int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userID int64) error
	List(ctx context.Context, limit, offset int) ([]models.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IuserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if user.UserInfo != nil {
			user.UserInfo.UserID = user.ID
			if err := tx.Create(user.UserInfo).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepository) GetByID(ctx context.Context, userID int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByIDWithInfo(ctx context.Context, userID int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("UserInfo").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserInfo{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.User{}, "id = ?", userID).Error
	})
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	var list []models.User
	query := r.db.WithContext(ctx).Model(&models.User{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	return list, total, err
}
