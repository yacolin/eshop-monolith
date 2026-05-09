package repositories

import (
	"context"
	userModels "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/infra/repository/models"
	"time"

	"gorm.io/gorm"
)

// AuthTokenRepository 认证令牌仓库接口
type IauthTokenRepository interface {
	// Create 创建令牌记录
	Create(ctx context.Context, token *userModels.AuthToken) error
	// GetByID 根据ID获取令牌
	GetByID(ctx context.Context, id string) (*userModels.AuthToken, error)
	// GetByJTI 根据JTI获取令牌
	GetByJTI(ctx context.Context, jti string) (*userModels.AuthToken, error)
	// GetByUserID 获取用户的所有令牌
	GetByUserID(ctx context.Context, userID string) ([]userModels.AuthToken, error)
	// GetActiveByUserID 获取用户的有效令牌
	GetActiveByUserID(ctx context.Context, userID string) ([]userModels.AuthToken, error)
	// Revoke 撤销令牌
	Revoke(ctx context.Context, jti string) error
	// RevokeAllByUserID 撤销用户的所有令牌
	RevokeAllByUserID(ctx context.Context, userID string) error
	// IsRevoked 检查令牌是否已撤销
	IsRevoked(ctx context.Context, jti string) (bool, error)
	// DeleteExpired 删除过期的令牌
	DeleteExpired(ctx context.Context, before time.Time) error
	// Update 更新令牌
	Update(ctx context.Context, token *userModels.AuthToken) error
	// Delete 删除令牌
	Delete(ctx context.Context, id string) error
}

// LoginHistoryRepository 登录历史仓库接口
type IloginHistoryRepository interface {
	// Create 创建登录历史记录
	Create(ctx context.Context, history *userModels.LoginHistory) error
	// GetByUserID 获取用户的登录历史
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]userModels.LoginHistory, int64, error)
	// GetByID 根据ID获取登录历史
	GetByID(ctx context.Context, id string) (*userModels.LoginHistory, error)
}

type authTokenRepository struct {
	db *gorm.DB
}

// NewAuthTokenRepository 创建认证令牌仓库实例
func NewAuthTokenRepository(db *gorm.DB) IauthTokenRepository {
	return &authTokenRepository{db: db}
}

func (r *authTokenRepository) Create(ctx context.Context, token *userModels.AuthToken) error {
	po := models.AuthTokenFromDomain(token)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	token.ID = po.ID
	return nil
}

func (r *authTokenRepository) GetByID(ctx context.Context, id string) (*userModels.AuthToken, error) {
	var po models.AuthTokenPO
	err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *authTokenRepository) GetByJTI(ctx context.Context, jti string) (*userModels.AuthToken, error) {
	var po models.AuthTokenPO
	err := r.db.WithContext(ctx).Where("jti = ?", jti).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *authTokenRepository) GetByUserID(ctx context.Context, userID string) ([]userModels.AuthToken, error) {
	var pos []models.AuthTokenPO
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}

	tokens := make([]userModels.AuthToken, len(pos))
	for i, po := range pos {
		tokens[i] = *po.ToDomain()
	}
	return tokens, nil
}

func (r *authTokenRepository) GetActiveByUserID(ctx context.Context, userID string) ([]userModels.AuthToken, error) {
	pos := make([]models.AuthTokenPO, 0)
	err := r.db.WithContext(ctx).Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}

	tokens := make([]userModels.AuthToken, len(pos))
	for i, po := range pos {
		tokens[i] = *po.ToDomain()
	}
	return tokens, nil
}

func (r *authTokenRepository) Revoke(ctx context.Context, jti string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.AuthTokenPO{}).Where("jti = ?", jti).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
			"updated_at": now,
		}).Error
}

func (r *authTokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.AuthTokenPO{}).Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
			"updated_at": now,
		}).Error
}

func (r *authTokenRepository) IsRevoked(ctx context.Context, jti string) (bool, error) {
	var po models.AuthTokenPO
	err := r.db.WithContext(ctx).Where("jti = ?", jti).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil // 记录不存在，视为已撤销
		}
		return false, err
	}
	return po.Revoked || po.ExpiresAt.Before(time.Now()), nil
}

func (r *authTokenRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", before).Delete(&models.AuthTokenPO{}).Error
}

func (r *authTokenRepository) Update(ctx context.Context, token *userModels.AuthToken) error {
	po := models.AuthTokenFromDomain(token)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *authTokenRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AuthTokenPO{}, "id = ?", id).Error
}

// LoginHistory Repository 实现

type loginHistoryRepository struct {
	db *gorm.DB
}

// NewLoginHistoryRepository 创建登录历史仓库实例
func NewLoginHistoryRepository(db *gorm.DB) IloginHistoryRepository {
	return &loginHistoryRepository{db: db}
}

func (r *loginHistoryRepository) Create(ctx context.Context, history *userModels.LoginHistory) error {
	po := models.LoginHistoryFromDomain(history)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	history.ID = po.ID
	return nil
}

func (r *loginHistoryRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]userModels.LoginHistory, int64, error) {
	var pos []models.LoginHistoryPO
	var total int64

	query := r.db.WithContext(ctx).Model(&models.LoginHistoryPO{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	histories := make([]userModels.LoginHistory, len(pos))
	for i, po := range pos {
		histories[i] = *po.ToDomain()
	}
	return histories, total, nil
}

func (r *loginHistoryRepository) GetByID(ctx context.Context, id string) (*userModels.LoginHistory, error) {
	var po models.LoginHistoryPO
	err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}
