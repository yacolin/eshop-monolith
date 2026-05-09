package repositories

import (
	"context"
	userModels "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

type IuserInfoRepository interface {
	CreateUserInfo(ctx context.Context, userInfo *userModels.UserInfo) error
	GetUserInfoByUserID(ctx context.Context, userID int64) (*userModels.UserInfo, error)
	UpdateUserInfo(ctx context.Context, userInfo *userModels.UserInfo) error
}

type userInfoRepository struct {
	db *gorm.DB
}

func NewUserInfoRepository(db *gorm.DB) IuserInfoRepository {
	return &userInfoRepository{db: db}
}

func (r *userInfoRepository) CreateUserInfo(ctx context.Context, userInfo *userModels.UserInfo) error {
	po := models.UserInfoFromDomain(userInfo)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	userInfo.ID = po.ID
	return nil
}

func (r *userInfoRepository) GetUserInfoByUserID(ctx context.Context, userID int64) (*userModels.UserInfo, error) {
	var po models.UserInfoPO
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *userInfoRepository) UpdateUserInfo(ctx context.Context, userInfo *userModels.UserInfo) error {
	po := models.UserInfoFromDomain(userInfo)
	return r.db.WithContext(ctx).Save(po).Error
}
