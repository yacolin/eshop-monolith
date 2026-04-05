package repositories

import (
	"context"
	"eshop-monolith/internal/user/domain/models"

	"gorm.io/gorm"
)

type IuserIdentityRepository interface {
	// Create 创建身份凭证
	Create(ctx context.Context, identity *models.UserIdentity) error
	// GetByID 根据ID获取身份凭证
	GetByID(ctx context.Context, id int64) (*models.UserIdentity, error)
	// GetByProviderAndIdentifier 根据provider和identifier获取身份凭证
	GetByProviderAndIdentifier(ctx context.Context, provider, identifier string) (*models.UserIdentity, error)
	// GetByUserID 根据用户ID获取所有身份凭证
	GetByUserID(ctx context.Context, userID string) ([]models.UserIdentity, error)
	// GetByUserIDAndProvider 根据用户ID和provider获取特定类型的身份凭证
	GetByUserIDAndProvider(ctx context.Context, userID, provider string) (*models.UserIdentity, error)
	// Update 更新身份凭证
	Update(ctx context.Context, identity *models.UserIdentity) error
	// Delete 删除身份凭证
	Delete(ctx context.Context, id string) error
	// DeleteByUserID 删除用户的所有身份凭证
	DeleteByUserID(ctx context.Context, userID string) error
	// Exists 检查身份凭证是否存在
	Exists(ctx context.Context, provider, identifier string) (bool, error)
	// LinkIdentityToUser 将身份凭证关联到用户
	LinkIdentityToUser(ctx context.Context, identityID, userID string) error
	// GetUserByIdentity 根据身份凭证获取用户信息
	GetUserByIdentity(ctx context.Context, provider, identifier string) (*models.User, error)
}

type userIdentityRepository struct {
	db *gorm.DB
}

// NewUserIdentityRepository 创建用户身份凭证仓库实例
func NewUserIdentityRepository(db *gorm.DB) IuserIdentityRepository {
	return &userIdentityRepository{db: db}
}

func (r *userIdentityRepository) Create(ctx context.Context, identity *models.UserIdentity) error {
	return r.db.WithContext(ctx).Create(identity).Error
}

func (r *userIdentityRepository) GetByID(ctx context.Context, id int64) (*models.UserIdentity, error) {
	var identity models.UserIdentity
	err := r.db.WithContext(ctx).Preload("User").First(&identity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *userIdentityRepository) GetByProviderAndIdentifier(ctx context.Context, provider, identifier string) (*models.UserIdentity, error) {
	var identity models.UserIdentity
	err := r.db.WithContext(ctx).Preload("User").Where("provider = ? AND identifier = ?", provider, identifier).First(&identity).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *userIdentityRepository) GetByUserID(ctx context.Context, userID string) ([]models.UserIdentity, error) {
	var identities []models.UserIdentity
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&identities).Error
	if err != nil {
		return nil, err
	}
	return identities, nil
}

func (r *userIdentityRepository) GetByUserIDAndProvider(ctx context.Context, userID, provider string) (*models.UserIdentity, error) {
	var identity models.UserIdentity
	err := r.db.WithContext(ctx).Where("user_id = ? AND provider = ?", userID, provider).First(&identity).Error
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *userIdentityRepository) Update(ctx context.Context, identity *models.UserIdentity) error {
	return r.db.WithContext(ctx).Save(identity).Error
}

func (r *userIdentityRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.UserIdentity{}, "id = ?", id).Error
}

func (r *userIdentityRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&models.UserIdentity{}, "user_id = ?", userID).Error
}

func (r *userIdentityRepository) Exists(ctx context.Context, provider, identifier string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.UserIdentity{}).Where("provider = ? AND identifier = ?", provider, identifier).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userIdentityRepository) LinkIdentityToUser(ctx context.Context, identityID, userID string) error {
	return r.db.WithContext(ctx).Model(&models.UserIdentity{}).Where("id = ?", identityID).Update("user_id", userID).Error
}

func (r *userIdentityRepository) GetUserByIdentity(ctx context.Context, provider, identifier string) (*models.User, error) {
	var identity models.UserIdentity
	err := r.db.WithContext(ctx).Preload("User").Where("provider = ? AND identifier = ?", provider, identifier).First(&identity).Error
	if err != nil {
		return nil, err
	}
	if identity.User == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return identity.User, nil
}
