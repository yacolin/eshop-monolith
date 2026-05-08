package service

import (
	"context"
	"time"

	"eshop-monolith/internal/eventbus"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/internal/user/api/dto"
	"eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/user/domain/repositories"
	"eshop-monolith/internal/user/events"

	"gorm.io/gorm"
)

type UserService struct {
	repo             repositories.IuserRepository
	userInfoRepo     repositories.IuserInfoRepository
	authRepo         repositories.IauthTokenRepository
	loginHistoryRepo repositories.IloginHistoryRepository
	bus              *eventbus.Bus

	jwtSecret string
}

func NewUserService(repo repositories.IuserRepository, userInfoRepo repositories.IuserInfoRepository, authRepo repositories.IauthTokenRepository, loginHistoryRepo repositories.IloginHistoryRepository, bus *eventbus.Bus) *UserService {
	return &UserService{
		repo:             repo,
		userInfoRepo:     userInfoRepo,
		authRepo:         authRepo,
		loginHistoryRepo: loginHistoryRepo,
		bus:              bus,
	}
}

func (s *UserService) SetJWTSecret(secret string) {
	s.jwtSecret = secret
}

// GetProfile 获取用户资料（包含 UserInfo）
func (s *UserService) GetProfile(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.repo.GetByIDWithInfo(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetUserInfo 获取用户详细信息
func (s *UserService) GetUserInfo(ctx context.Context, userID int64) (*models.UserInfo, error) {
	userInfo, err := s.userInfoRepo.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return userInfo, nil
}

// UpdateUserInfo 更新用户详细信息（包含 Avatar、Nickname 等）
func (s *UserService) UpdateUserInfo(ctx context.Context, userID int64, req dto.UpdateUserInfoRequest) (*models.UserInfo, error) {
	userInfo, err := s.userInfoRepo.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userInfo = &models.UserInfo{UserID: userID}
		} else {
			return nil, err
		}
	}

	if req.Nickname != "" {
		userInfo.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		userInfo.Avatar = req.Avatar
	}
	if req.Gender != 0 {
		userInfo.Gender = req.Gender
	}
	if req.Birthday != "" {
		// Parse birthday string to time.Time
		// Assuming format: 2006-01-02
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			userInfo.Birthday = &birthday
		}
	}
	if req.Address != "" {
		userInfo.Address = req.Address
	}
	if req.Bio != "" {
		userInfo.Bio = req.Bio
	}
	if req.Country != "" {
		userInfo.Country = req.Country
	}
	if req.Province != "" {
		userInfo.Province = req.Province
	}
	if req.City != "" {
		userInfo.City = req.City
	}
	if req.ZipCode != "" {
		userInfo.ZipCode = req.ZipCode
	}
	if req.Language != "" {
		userInfo.Language = req.Language
	}
	if req.Timezone != "" {
		userInfo.Timezone = req.Timezone
	}

	if userInfo.ID == 0 {
		if err := s.userInfoRepo.CreateUserInfo(ctx, userInfo); err != nil {
			return nil, err
		}
	} else {
		if err := s.userInfoRepo.UpdateUserInfo(ctx, userInfo); err != nil {
			return nil, err
		}
	}

	s.bus.Publish(events.UserProfileUpdatedEvent{
		UserID:   userID,
		Nickname: userInfo.Nickname,
	})

	return userInfo, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}
