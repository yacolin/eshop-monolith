package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/token"
	"eshop-monolith/pkg/utils"
)

type AuthService struct {
	db               *gorm.DB
	userRepo         IuserRepository
	infoRepo         IuserInfoRepository
	loginHistoryRepo IloginHistoryRepository
}

func NewAuthService(db *gorm.DB, userRepo IuserRepository, infoRepo IuserInfoRepository, loginHistoryRepo IloginHistoryRepository) *AuthService {
	return &AuthService{
		db:               db,
		userRepo:         userRepo,
		infoRepo:         infoRepo,
		loginHistoryRepo: loginHistoryRepo,
	}
}

func (s *AuthService) LoginByPassword(ctx context.Context, username, password string) (*User, *token.TokenPair, error) {
	if username == "" || password == "" {
		return nil, nil, errcode.ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errcode.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	if user.Status != 1 {
		return nil, nil, errcode.ErrAccountDisabled
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, nil, errcode.ErrInvalidCredentials
	}

	// C 端用户已无角色体系(usr_roles 等表已从库中删除),roles 为空
	tokenPair, err := token.GenerateTokenPair(user.ID, user.Username, nil)
	if err != nil {
		return nil, nil, err
	}

	s.RecordLoginHistory(ctx, user.ID, "password", LoginEventLogin, LoginStatusSuccess, "", "", "")
	return user, tokenPair, nil
}

func (s *AuthService) Register(ctx context.Context, req *RegisterReq) (*User, *token.TokenPair, error) {
	if req.Username != "" {
		_, err := s.userRepo.FindByUsername(ctx, req.Username)
		if err == nil {
			return nil, nil, errcode.ErrUsernameAlreadyExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
	}

	var user *User
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}

		user = &User{
			Username:     req.Username,
			PasswordHash: string(hash),
			Email:        req.Email,
			Phone:        req.Phone,
			Status:       1,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		info := &UserInfo{UserID: user.ID}
		return tx.Create(info).Error // 不再插入 usr_user_roles(表已删除)
	})
	if err != nil {
		return nil, nil, err
	}

	tokenPair, err := token.GenerateTokenPair(user.ID, user.Username, nil)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

func (s *AuthService) RecordLoginHistory(ctx context.Context, userID int64, provider, event, status, failReason, ip, userAgent string) {
	go func() {
		_ = s.loginHistoryRepo.Create(context.Background(), &LoginHistory{
			UserID:     userID,
			Provider:   provider,
			Event:      event,
			Status:     status,
			FailReason: failReason,
			IP:         ip,
			UserAgent:  userAgent,
		})
	}()
}
