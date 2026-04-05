package repositories

import (
	"context"
	"eshop-monolith/internal/user/domain/models"

	"gorm.io/gorm"
)

type IuserInfoRepository interface {
	CreateUserInfo(ctx context.Context, userInfo *models.UserInfo) error
	GetUserInfoByUserID(ctx context.Context, userID int64) (*models.UserInfo, error)
	UpdateUserInfo(ctx context.Context, userInfo *models.UserInfo) error
}

type userInfoRepository struct {
	db *gorm.DB
}

func NewUserInfoRepository(db *gorm.DB) IuserInfoRepository {
	return &userInfoRepository{db: db}
}

func (r *userInfoRepository) CreateUserInfo(ctx context.Context, userInfo *models.UserInfo) error {
	return r.db.WithContext(ctx).Create(userInfo).Error
}

func (r *userInfoRepository) GetUserInfoByUserID(ctx context.Context, userID int64) (*models.UserInfo, error) {
	var userInfo models.UserInfo
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (r *userInfoRepository) UpdateUserInfo(ctx context.Context, userInfo *models.UserInfo) error {
	return r.db.WithContext(ctx).Save(userInfo).Error
}
