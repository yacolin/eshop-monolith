package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/utils"
	"eshop-monolith/internal/user/domain/auth"
	"eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/user/domain/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db                *gorm.DB
	userRepo          repositories.IuserRepository
	identityRepo      repositories.IuserIdentityRepository
	tokenRepo         repositories.IauthTokenRepository
	loginHistoryRepo  repositories.IloginHistoryRepository
	tokenService      *TokenService
	wechatClient      WechatClient
	verifyCodeService VerifyCodeService
}

// WechatClient 微信客户端接口
type WechatClient interface {
	Code2Session(ctx context.Context, appID, secret, code string) (*auth.WechatSession, error)
}

// VerifyCodeService 验证码服务接口
type VerifyCodeService interface {
	Verify(ctx context.Context, phoneOrEmail, code string) error
	Send(ctx context.Context, phoneOrEmail string) error
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	db *gorm.DB,
	userRepo repositories.IuserRepository,
	identityRepo repositories.IuserIdentityRepository,
	tokenRepo repositories.IauthTokenRepository,
	loginHistoryRepo repositories.IloginHistoryRepository,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		db:               db,
		userRepo:         userRepo,
		identityRepo:     identityRepo,
		tokenRepo:        tokenRepo,
		loginHistoryRepo: loginHistoryRepo,
		tokenService:     tokenService,
	}
}

// SetWechatClient 设置微信客户端
func (s *AuthService) SetWechatClient(client WechatClient) {
	s.wechatClient = client
}

// SetVerifyCodeService 设置验证码服务
func (s *AuthService) SetVerifyCodeService(svc VerifyCodeService) {
	s.verifyCodeService = svc
}

// LoginByPassword 用户名密码登录
func (s *AuthService) LoginByPassword(ctx context.Context, payload *auth.PasswordPayload) (*models.User, *models.UserIdentity, error) {
	// 1. 查询用户身份凭证
	identity, err := s.identityRepo.GetByProviderAndIdentifier(ctx, models.ProviderPassword.String(), payload.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, errcode.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(identity.Credential), []byte(payload.Password)); err != nil {
		// 密码验证失败，返回无效凭证错误
		return nil, nil, errcode.ErrInvalidCredentials
	}

	// 3. 获取用户信息
	user, err := s.userRepo.GetByID(ctx, identity.UserID)
	if err != nil {
		return nil, nil, err
	}

	// 4. 检查用户状态
	if user.Status != 1 {
		return nil, nil, errcode.ErrAccountDisabled
	}

	return user, identity, nil
}

// LoginByWechat 微信登录
func (s *AuthService) LoginByWechat(ctx context.Context, payload *auth.WechatPayload, appSecret string) (*models.User, *models.UserIdentity, bool, error) {
	// 1. 调用微信接口获取openid和session_key
	if s.wechatClient == nil {
		return nil, nil, false, errcode.ErrWechatClientNotConfigured
	}

	session, err := s.wechatClient.Code2Session(ctx, payload.AppID, appSecret, payload.Code)
	if err != nil {
		return nil, nil, false, fmt.Errorf("微信登录失败: %w", err)
	}

	if session.ErrCode != 0 {
		return nil, nil, false, fmt.Errorf("微信登录失败: %s", session.ErrMsg)
	}

	// 2. 查询是否已存在该微信身份
	identity, err := s.identityRepo.GetByProviderAndIdentifier(ctx, models.ProviderWechat.String(), session.OpenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false, err
	}

	// 3. 如果存在，返回用户信息
	if identity != nil {
		// 更新session_key（可选，建议加密存储）
		meta := models.IdentityMeta{
			UnionID:    session.UnionID,
			SessionKey: session.SessionKey, // 实际生产环境应该加密存储
			Source:     payload.Source,
			AppID:      payload.AppID,
		}
		metaJSON, _ := json.Marshal(meta)
		identity.Meta = string(metaJSON)
		_ = s.identityRepo.Update(ctx, identity)

		user, err := s.userRepo.GetByID(ctx, identity.UserID)
		if err != nil {
			return nil, nil, false, err
		}
		return user, identity, false, nil
	}

	// 4. 如果不存在，创建新的身份凭证（但不创建用户，需要后续绑定或注册）
	meta := models.IdentityMeta{
		UnionID:    session.UnionID,
		SessionKey: session.SessionKey,
		Source:     payload.Source,
		AppID:      payload.AppID,
	}
	metaJSON, _ := json.Marshal(meta)

	newIdentity := &models.UserIdentity{
		Provider:   models.ProviderWechat.String(),
		Identifier: session.OpenID,
		Verified:   true, // 微信登录视为已验证
		Meta:       string(metaJSON),
	}

	// 注意：此时还没有关联UserID，需要后续绑定
	return nil, newIdentity, true, nil
}

// LoginByPhone 手机号验证码登录
func (s *AuthService) LoginByPhone(ctx context.Context, payload *auth.PhonePayload) (*models.User, *models.UserIdentity, bool, error) {
	// 1. 验证验证码
	if s.verifyCodeService != nil {
		if err := s.verifyCodeService.Verify(ctx, payload.Phone, payload.VerifyCode); err != nil {
			return nil, nil, false, fmt.Errorf("验证码错误: %w", err)
		}
	}

	// 2. 查询是否已存在该手机号身份
	identity, err := s.identityRepo.GetByProviderAndIdentifier(ctx, models.ProviderPhone.String(), payload.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false, err
	}

	// 3. 如果存在，返回用户信息
	if identity != nil {
		user, err := s.userRepo.GetByID(ctx, identity.UserID)
		if err != nil {
			return nil, nil, false, err
		}
		return user, identity, false, nil
	}

	// 4. 如果不存在，创建新的身份凭证
	newIdentity := &models.UserIdentity{
		Provider:   models.ProviderPhone.String(),
		Identifier: payload.Phone,
		Verified:   true,
	}

	return nil, newIdentity, true, nil
}

// LoginByEmail 邮箱验证码登录
func (s *AuthService) LoginByEmail(ctx context.Context, payload *auth.EmailPayload) (*models.User, *models.UserIdentity, bool, error) {
	// 1. 验证验证码
	if s.verifyCodeService != nil {
		if err := s.verifyCodeService.Verify(ctx, payload.Email, payload.VerifyCode); err != nil {
			return nil, nil, false, fmt.Errorf("验证码错误: %w", err)
		}
	}

	// 2. 查询是否已存在该邮箱身份
	identity, err := s.identityRepo.GetByProviderAndIdentifier(ctx, models.ProviderEmail.String(), payload.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false, err
	}

	// 3. 如果存在，返回用户信息
	if identity != nil {
		user, err := s.userRepo.GetByID(ctx, identity.UserID)
		if err != nil {
			return nil, nil, false, err
		}
		return user, identity, false, nil
	}

	// 4. 如果不存在，创建新的身份凭证
	newIdentity := &models.UserIdentity{
		Provider:   models.ProviderEmail.String(),
		Identifier: payload.Email,
		Verified:   true,
	}

	return nil, newIdentity, true, nil
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, payload *auth.RegisterPayload) (*models.User, *models.UserIdentity, error) {
	// 1. 检查用户名是否已存在
	if payload.Username != "" {
		exists, err := s.identityRepo.Exists(ctx, models.ProviderPassword.String(), payload.Username)
		if err != nil {
			return nil, nil, err
		}
		if exists {
			return nil, nil, errcode.ErrUsernameAlreadyExists
		}
	}

	// 2. 开始事务
	var user *models.User
	var identity *models.UserIdentity

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建用户 - 只保留业务字段
		user = &models.User{
			Status: 1,
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 创建用户详情
		userInfo := &models.UserInfo{
			UserID: user.ID,
		}
		if err := tx.Create(userInfo).Error; err != nil {
			return err
		}

		// 创建身份凭证
		switch payload.Provider {
		case models.ProviderPassword.String():
			// 密码加密
			hashedPassword, err := utils.HashPassword(payload.Password)
			if err != nil {
				return err
			}
			identity = &models.UserIdentity{
				UserID:     user.ID,
				Provider:   models.ProviderPassword.String(),
				Identifier: payload.Username,
				Credential: string(hashedPassword),
				Verified:   true,
				Meta:       "{}",
			}
		case models.ProviderPhone.String():
			identity = &models.UserIdentity{
				UserID:     user.ID,
				Provider:   models.ProviderPhone.String(),
				Identifier: payload.Phone,
				Verified:   true,
				Meta:       "{}",
			}
		case models.ProviderEmail.String():
			identity = &models.UserIdentity{
				UserID:     user.ID,
				Provider:   models.ProviderEmail.String(),
				Identifier: payload.Email,
				Verified:   true,
				Meta:       "{}",
			}
		default:
			return errcode.ErrUnsupportedProvider
		}

		if err := tx.Create(identity).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return user, identity, nil
}

// BindIdentity 绑定身份凭证到用户
func (s *AuthService) BindIdentity(ctx context.Context, userID string, identity *models.UserIdentity) error {
	// 1. 检查用户是否存在
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}
	_, err = s.userRepo.GetByID(ctx, userIDInt)
	if err != nil {
		return err
	}

	// 2. 检查身份凭证是否已存在
	exists, err := s.identityRepo.Exists(ctx, identity.Provider, identity.Identifier)
	if err != nil {
		return err
	}
	if exists {
		return errcode.ErrIdentityAlreadyBound
	}

	// 3. 绑定身份凭证
	identity.UserID = userIDInt
	return s.identityRepo.Create(ctx, identity)
}

// UnbindIdentity 解绑身份凭证
func (s *AuthService) UnbindIdentity(ctx context.Context, userID, provider string) error {
	identity, err := s.identityRepo.GetByUserIDAndProvider(ctx, userID, provider)
	if err != nil {
		return err
	}

	return s.identityRepo.Delete(ctx, strconv.FormatInt(identity.ID, 10))
}

// GetUserIdentities 获取用户的所有身份凭证
func (s *AuthService) GetUserIdentities(ctx context.Context, userID string) ([]models.UserIdentity, error) {
	return s.identityRepo.GetByUserID(ctx, userID)
}

// RecordLoginHistory 记录登录历史
func (s *AuthService) RecordLoginHistory(ctx context.Context, userID int64, identityID int64, provider, event, status, failReason, ip, userAgent string) {
	if s.loginHistoryRepo == nil {
		return
	}

	history := &models.LoginHistory{
		UserID:     userID,
		IdentityID: identityID,
		Provider:   provider,
		IP:         ip,
		UserAgent:  userAgent,
		Event:      event,
		Status:     status,
		FailReason: failReason,
	}

	_ = s.loginHistoryRepo.Create(ctx, history)
}

// HashPassword 密码加密
func (s *AuthService) HashPassword(password string) (string, error) {
	return utils.HashPassword(password)
}

// VerifyPassword 验证密码
func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	if !utils.CheckPasswordHash(password, hashedPassword) {
		return errcode.ErrInvalidCredentials
	}
	return nil
}

// GenerateRandomUsername 生成随机用户名
func (s *AuthService) GenerateRandomUsername() string {
	return "user_" + utils.GenerateRandomString(8)
}
