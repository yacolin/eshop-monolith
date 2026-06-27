package repositories

import (
	"context"
	userModels "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

type IuserIdentityRepository interface {
	// Create 创建身份凭证
	Create(ctx context.Context, identity *userModels.UserIdentity) error
	// GetByID 根据ID获取身份凭证
	GetByID(ctx context.Context, id int64) (*userModels.UserIdentity, error)
	// GetByProviderAndIdentifier 根据provider和identifier获取身份凭证
	GetByProviderAndIdentifier(ctx context.Context, provider, identifier string) (*userModels.UserIdentity, error)
	// GetByUserID 根据用户ID获取所有身份凭证
	GetByUserID(ctx context.Context, userID string) ([]userModels.UserIdentity, error)
	// GetByUserIDAndProvider 根据用户ID和provider获取特定类型的身份凭证
	GetByUserIDAndProvider(ctx context.Context, userID, provider string) (*userModels.UserIdentity, error)
	// Update 更新身份凭证
	Update(ctx context.Context, identity *userModels.UserIdentity) error
	// Delete 删除身份凭证
	Delete(ctx context.Context, id string) error
	// DeleteByUserID 删除用户的所有身份凭证
	DeleteByUserID(ctx context.Context, userID string) error
	// Exists 检查身份凭证是否存在
	Exists(ctx context.Context, provider, identifier string) (bool, error)
	// LinkIdentityToUser 将身份凭证关联到用户
	LinkIdentityToUser(ctx context.Context, identityID, userID string) error
	// GetUserByIdentity 根据身份凭证获取用户信息
	GetUserByIdentity(ctx context.Context, provider, identifier string) (*userModels.User, error)
	// GetWithUserByProviderAndIdentifier 单次 JOIN 查询身份+用户（替代 Preload("User")）
	GetWithUserByProviderAndIdentifier(ctx context.Context, provider, identifier string) (*userModels.UserIdentity, error)
}

type userIdentityRepository struct {
	db *gorm.DB
}

// NewUserIdentityRepository 创建用户身份凭证仓库实例
func NewUserIdentityRepository(db *gorm.DB) IuserIdentityRepository {
	return &userIdentityRepository{db: db}
}

func (r *userIdentityRepository) Create(ctx context.Context, identity *userModels.UserIdentity) error {
	po := models.UserIdentityFromDomain(identity)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	identity.ID = po.ID
	return nil
}

func (r *userIdentityRepository) GetByID(ctx context.Context, id int64) (*userModels.UserIdentity, error) {
	var po models.UserIdentityPO
	err := r.db.WithContext(ctx).Preload("User").First(&po, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *userIdentityRepository) GetByProviderAndIdentifier(ctx context.Context, provider, identifier string) (*userModels.UserIdentity, error) {
	var po models.UserIdentityPO
	err := r.db.WithContext(ctx).Preload("User").Where("provider = ? AND identifier = ?", provider, identifier).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *userIdentityRepository) GetByUserID(ctx context.Context, userID string) ([]userModels.UserIdentity, error) {
	var pos []models.UserIdentityPO
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	identities := make([]userModels.UserIdentity, len(pos))
	for i, po := range pos {
		identities[i] = *po.ToDomain()
	}
	return identities, nil
}

func (r *userIdentityRepository) GetByUserIDAndProvider(ctx context.Context, userID, provider string) (*userModels.UserIdentity, error) {
	var po models.UserIdentityPO
	err := r.db.WithContext(ctx).Where("user_id = ? AND provider = ?", userID, provider).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *userIdentityRepository) Update(ctx context.Context, identity *userModels.UserIdentity) error {
	po := models.UserIdentityFromDomain(identity)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *userIdentityRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.UserIdentityPO{}, "id = ?", id).Error
}

func (r *userIdentityRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&models.UserIdentityPO{}, "user_id = ?", userID).Error
}

func (r *userIdentityRepository) Exists(ctx context.Context, provider, identifier string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.UserIdentityPO{}).Where("provider = ? AND identifier = ?", provider, identifier).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userIdentityRepository) LinkIdentityToUser(ctx context.Context, identityID, userID string) error {
	return r.db.WithContext(ctx).Model(&models.UserIdentityPO{}).Where("id = ?", identityID).Update("user_id", userID).Error
}

func (r *userIdentityRepository) GetUserByIdentity(ctx context.Context, provider, identifier string) (*userModels.User, error) {
	var po models.UserIdentityPO
	err := r.db.WithContext(ctx).Preload("User").Where("provider = ? AND identifier = ?", provider, identifier).First(&po).Error
	if err != nil {
		return nil, err
	}
	if po.User == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return po.User.ToDomain(), nil
}

// GetWithUserByProviderAndIdentifier 单次 JOIN 查询身份+用户（替代 Preload("User")，减少1次DB往返）
func (r *userIdentityRepository) GetWithUserByProviderAndIdentifier(ctx context.Context, provider, identifier string) (*userModels.UserIdentity, error) {
	var po models.UserIdentityPO
	err := r.db.WithContext(ctx).
		Joins("User").
		Where("user_identities.provider = ? AND user_identities.identifier = ?", provider, identifier).
		First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}
