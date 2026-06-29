package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/utils"
)

type AuthService struct {
	db               *gorm.DB
	userRepo         IuserRepository
	infoRepo         IuserInfoRepository
	roleRepo         IroleRepository
	loginHistoryRepo IloginHistoryRepository
	tokenSvc         *TokenService
}

func NewAuthService(db *gorm.DB, userRepo IuserRepository, infoRepo IuserInfoRepository, roleRepo IroleRepository, loginHistoryRepo IloginHistoryRepository, tokenSvc *TokenService) *AuthService {
	return &AuthService{
		db:               db,
		userRepo:         userRepo,
		infoRepo:         infoRepo,
		roleRepo:         roleRepo,
		loginHistoryRepo: loginHistoryRepo,
		tokenSvc:         tokenSvc,
	}
}

func (s *AuthService) LoginByPassword(ctx context.Context, username, password string) (*User, *TokenPair, error) {
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

	tokenPair, err := s.tokenSvc.GenerateTokenPair(ctx, user.ID, user.Username)
	if err != nil {
		return nil, nil, err
	}

	s.RecordLoginHistory(ctx, user.ID, "password", LoginEventLogin, LoginStatusSuccess, "", "", "")
	return user, tokenPair, nil
}

func (s *AuthService) Register(ctx context.Context, req *RegisterReq) (*User, *TokenPair, error) {
	if req.Username != "" {
		existing, err := s.userRepo.FindByUsername(ctx, req.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		if existing != nil {
			return nil, nil, errcode.ErrUsernameAlreadyExists
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
		if err := tx.Create(info).Error; err != nil {
			return err
		}

		var roleID int64
		if err := tx.Raw("SELECT id FROM usr_roles WHERE name = ?", "user").Scan(&roleID).Error; err != nil {
			return err
		}
		return tx.Exec("INSERT INTO usr_user_roles (user_id, role_id) VALUES (?, ?)", user.ID, roleID).Error
	})
	if err != nil {
		return nil, nil, err
	}

	tokenPair, err := s.tokenSvc.GenerateTokenPair(ctx, user.ID, user.Username)
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
