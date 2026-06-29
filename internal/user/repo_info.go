package user

import (
	"context"

	"gorm.io/gorm"
)

type IuserInfoRepository interface {
	CreateUserInfo(ctx context.Context, info *UserInfo) error
	GetUserInfoByUserID(ctx context.Context, userID int64) (*UserInfo, error)
	UpdateUserInfo(ctx context.Context, info *UserInfo) error
}

type UserInfoRepository struct {
	db *gorm.DB
}

func NewUserInfoRepository(db *gorm.DB) IuserInfoRepository {
	return &UserInfoRepository{db: db}
}

func (r *UserInfoRepository) CreateUserInfo(ctx context.Context, info *UserInfo) error {
	return r.db.WithContext(ctx).Create(info).Error
}

func (r *UserInfoRepository) GetUserInfoByUserID(ctx context.Context, userID int64) (*UserInfo, error) {
	var info UserInfo
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&info).Error
	return &info, err
}

func (r *UserInfoRepository) UpdateUserInfo(ctx context.Context, info *UserInfo) error {
	return r.db.WithContext(ctx).Model(info).Updates(info).Error
}
