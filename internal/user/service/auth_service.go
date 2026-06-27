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

type WechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// AuthService 认证服务
type AuthService struct {
	db                *gorm.DB
	identityRepo      repositories.IuserIdentityRepository
	tokenRepo         repositories.IauthTokenRepository
	loginHistoryRepo  repositories.IloginHistoryRepository
	roleRepo          repositories.IroleRepository
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
	_ repositories.IuserRepository, // 保留兼容，构造后不再持有
	identityRepo repositories.IuserIdentityRepository,
	tokenRepo repositories.IauthTokenRepository,
	loginHistoryRepo repositories.IloginHistoryRepository,
	roleRepo repositories.IroleRepository,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		db:               db,
		identityRepo:     identityRepo,
		tokenRepo:        tokenRepo,
		loginHistoryRepo: loginHistoryRepo,
		roleRepo:         roleRepo,
		tokenService:     tokenService,
	}
}

// SetWechatClient 设置微信客户端
func (s *AuthService) SetWechatClient(client WechatClient) {
	s.wechatClient = client
}

// SetVerifyCodeService 设置验证码服务
func (s *AuthService) SetVerifyCodeService(service VerifyCodeService) {
	s.verifyCodeService = service
}

// LoginByPassword 用户名密码登录
func (s *AuthService) LoginByPassword(ctx context.Context, payload *auth.PasswordPayload) (*models.User, *models.UserIdentity, error) {
	if payload.Username == "" || payload.Password == "" {
		return nil, nil, errcode.ErrInvalidCredentials
	}

	// GetByProviderAndIdentifier 已 Preload("User")，无需再单独查 user
	identity, err := s.identityRepo.GetWithUserByProviderAndIdentifier(ctx, models.ProviderPassword.String(), payload.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errcode.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identity.Credential), []byte(payload.Password)); err != nil {
		return nil, nil, errcode.ErrInvalidCredentials
	}

	return identity.User, identity, nil
}

// LoginByWechat 微信登录
func (s *AuthService) LoginByWechat(ctx context.Context, payload *auth.WechatPayload, appSecret string) (*models.User, *models.UserIdentity, bool, error) {
	if s.wechatClient == nil {
		return nil, nil, false, errcode.ErrWechatClientNotConfigured
	}

	session, err := s.wechatClient.Code2Session(ctx, payload.AppID, appSecret, payload.Code)
	if err != nil {
		return nil, nil, false, fmt.Errorf("微信登录失败: %w", err)
	}

	if session.OpenID == "" {
		return nil, nil, false, fmt.Errorf("微信登录失败: code无效")
	}

	identity, err := s.identityRepo.GetWithUserByProviderAndIdentifier(ctx, models.ProviderWechat.String(), session.OpenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false, err
	}

	if identity != nil {
		meta := models.IdentityMeta{
			UnionID:    session.UnionID,
			SessionKey: session.SessionKey,
			Source:     payload.Source,
			AppID:      payload.AppID,
		}
		metaJSON, _ := json.Marshal(meta)
		identity.Meta = string(metaJSON)
		_ = s.identityRepo.Update(ctx, identity)

		return identity.User, identity, false, nil
	}

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
		Verified:   true,
		Meta:       string(metaJSON),
	}

	return nil, newIdentity, true, nil
}

// LoginByPhone 手机号验证码登录
func (s *AuthService) LoginByPhone(ctx context.Context, payload *auth.PhonePayload) (*models.User, *models.UserIdentity, bool, error) {
	if s.verifyCodeService != nil {
		if err := s.verifyCodeService.Verify(ctx, payload.Phone, payload.VerifyCode); err != nil {
			return nil, nil, false, fmt.Errorf("验证码错误: %w", err)
		}
	}

	identity, err := s.identityRepo.GetWithUserByProviderAndIdentifier(ctx, models.ProviderPhone.String(), payload.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false, err
	}

	if identity != nil {
		return identity.User, identity, false, nil
	}

	newIdentity := &models.UserIdentity{
		Provider:   models.ProviderPhone.String(),
		Identifier: payload.Phone,
		Verified:   true,
	}

	return nil, newIdentity, true, nil
}

// LoginByEmail 邮箱验证码登录
func (s *AuthService) LoginByEmail(ctx context.Context, payload *auth.EmailPayload) (*models.User, *models.UserIdentity, bool, error) {
	identity, err := s.identityRepo.GetWithUserByProviderAndIdentifier(ctx, models.ProviderEmail.String(), payload.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false, err
	}

	if identity != nil {
		return identity.User, identity, false, nil
	}

	newIdentity := &models.UserIdentity{
		Provider:   models.ProviderEmail.String(),
		Identifier: payload.Email,
		Verified:   true,
	}

	return nil, newIdentity, true, nil
}

// Register 注册
func (s *AuthService) Register(ctx context.Context, payload *auth.RegisterPayload) (*models.User, *models.UserIdentity, error) {
	if payload.Username != "" {
		exists, err := s.identityRepo.Exists(ctx, models.ProviderPassword.String(), payload.Username)
		if err != nil {
			return nil, nil, err
		}
		if exists {
			return nil, nil, errcode.ErrUsernameAlreadyExists
		}
	}

	var user *models.User
	var identity *models.UserIdentity

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user = &models.User{Status: 1}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		userInfo := &models.UserInfo{UserID: user.ID}
		if err := tx.Create(userInfo).Error; err != nil {
			return err
		}

		switch payload.Provider {
		case models.ProviderPassword.String():
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

		// 分配默认角色
		var roleID int64
		if err := tx.Raw("SELECT id FROM roles WHERE name = ?", "user").Scan(&roleID).Error; err != nil {
			return err
		}
		_ = &models.UserRole{
			UserID: user.ID,
			RoleID: roleID,
		}
		if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", user.ID, roleID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return user, identity, nil
}

// RecordLoginHistory 记录登录历史（异步写入，不阻塞登录主流程）
func (s *AuthService) RecordLoginHistory(ctx context.Context, userID, identityID int64, provider, action, status, failureReason, ip, userAgent string) {
	go func() {
		loginHistory := &models.LoginHistory{
			UserID:     userID,
			IdentityID: identityID,
			Provider:   provider,
			Event:      action,
			Status:     status,
			FailReason: failureReason,
			IP:         ip,
			UserAgent:  userAgent,
		}
		_ = s.loginHistoryRepo.Create(context.Background(), loginHistory)
	}()
}

// createPasswordIdentity 创建密码身份凭证
func (s *AuthService) createPasswordIdentity(userID int64, payload *auth.RegisterPayload) (*models.UserIdentity, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	identity := &models.UserIdentity{
		UserID:     userID,
		Provider:   models.ProviderPassword.String(),
		Identifier: payload.Username,
		Credential: string(hashedPassword),
		Verified:   true,
		Meta:       "{}",
	}

	if err := s.identityRepo.Create(context.Background(), identity); err != nil {
		return nil, err
	}

	return identity, nil
}

// EnsureUser 确保用户存在（社交登录时使用）
func (s *AuthService) EnsureUser(ctx context.Context, provider, identifier string) (*models.User, *models.UserIdentity, bool, error) {
	identity, err := s.identityRepo.GetWithUserByProviderAndIdentifier(ctx, provider, identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, false, nil
		}
		return nil, nil, false, err
	}

	return identity.User, identity, false, nil
}

// getStringValue 从 JWT claims 中获取字符串值
func getStringValue(claims map[string]interface{}, key string) string {
	if v, ok := claims[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			return strconv.FormatInt(int64(val), 10)
		default:
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}
