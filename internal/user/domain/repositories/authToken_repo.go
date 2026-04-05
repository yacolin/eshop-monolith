package repositories

import (
	"context"
	"eshop-monolith/internal/user/domain/models"
	"time"

	"gorm.io/gorm"
)

// AuthTokenRepository 认证令牌仓库接口
type IauthTokenRepository interface {
	// Create 创建令牌记录
	Create(ctx context.Context, token *models.AuthToken) error
	// GetByID 根据ID获取令牌
	GetByID(ctx context.Context, id string) (*models.AuthToken, error)
	// GetByJTI 根据JTI获取令牌
	GetByJTI(ctx context.Context, jti string) (*models.AuthToken, error)
	// GetByUserID 获取用户的所有令牌
	GetByUserID(ctx context.Context, userID string) ([]models.AuthToken, error)
	// GetActiveByUserID 获取用户的有效令牌
	GetActiveByUserID(ctx context.Context, userID string) ([]models.AuthToken, error)
	// Revoke 撤销令牌
	Revoke(ctx context.Context, jti string) error
	// RevokeAllByUserID 撤销用户的所有令牌
	RevokeAllByUserID(ctx context.Context, userID string) error
	// IsRevoked 检查令牌是否已撤销
	IsRevoked(ctx context.Context, jti string) (bool, error)
	// DeleteExpired 删除过期的令牌
	DeleteExpired(ctx context.Context, before time.Time) error
	// Update 更新令牌
	Update(ctx context.Context, token *models.AuthToken) error
	// Delete 删除令牌
	Delete(ctx context.Context, id string) error
}

// LoginHistoryRepository 登录历史仓库接口
type IloginHistoryRepository interface {
	// Create 创建登录历史记录
	Create(ctx context.Context, history *models.LoginHistory) error
	// GetByUserID 获取用户的登录历史
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.LoginHistory, int64, error)
	// GetByID 根据ID获取登录历史
	GetByID(ctx context.Context, id string) (*models.LoginHistory, error)
}

type authTokenRepository struct {
	db *gorm.DB
}

// NewAuthTokenRepository 创建认证令牌仓库实例
func NewAuthTokenRepository(db *gorm.DB) IauthTokenRepository {
	return &authTokenRepository{db: db}
}

func (r *authTokenRepository) Create(ctx context.Context, token *models.AuthToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *authTokenRepository) GetByID(ctx context.Context, id string) (*models.AuthToken, error) {
	var token models.AuthToken
	err := r.db.WithContext(ctx).First(&token, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *authTokenRepository) GetByJTI(ctx context.Context, jti string) (*models.AuthToken, error) {
	var token models.AuthToken
	err := r.db.WithContext(ctx).Where("jti = ?", jti).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *authTokenRepository) GetByUserID(ctx context.Context, userID string) ([]models.AuthToken, error) {
	var tokens []models.AuthToken
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *authTokenRepository) GetActiveByUserID(ctx context.Context, userID string) ([]models.AuthToken, error) {
	var tokens []models.AuthToken
	err := r.db.WithContext(ctx).Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Order("created_at DESC").Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *authTokenRepository) Revoke(ctx context.Context, jti string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.AuthToken{}).Where("jti = ?", jti).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
			"updated_at": now,
		}).Error
}

func (r *authTokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.AuthToken{}).Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
			"updated_at": now,
		}).Error
}

func (r *authTokenRepository) IsRevoked(ctx context.Context, jti string) (bool, error) {
	var token models.AuthToken
	err := r.db.WithContext(ctx).Where("jti = ?", jti).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil // 记录不存在，视为已撤销
		}
		return false, err
	}
	return token.Revoked || token.ExpiresAt.Before(time.Now()), nil
}

func (r *authTokenRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", before).Delete(&models.AuthToken{}).Error
}

func (r *authTokenRepository) Update(ctx context.Context, token *models.AuthToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

func (r *authTokenRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AuthToken{}, "id = ?", id).Error
}

// LoginHistory Repository 实现

type loginHistoryRepository struct {
	db *gorm.DB
}

// NewLoginHistoryRepository 创建登录历史仓库实例
func NewLoginHistoryRepository(db *gorm.DB) IloginHistoryRepository {
	return &loginHistoryRepository{db: db}
}

func (r *loginHistoryRepository) Create(ctx context.Context, history *models.LoginHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *loginHistoryRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.LoginHistory, int64, error) {
	var histories []models.LoginHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&models.LoginHistory{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *loginHistoryRepository) GetByID(ctx context.Context, id string) (*models.LoginHistory, error) {
	var history models.LoginHistory
	err := r.db.WithContext(ctx).First(&history, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}
