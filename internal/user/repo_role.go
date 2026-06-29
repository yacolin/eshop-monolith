package user

import (
	"context"

	"gorm.io/gorm"
)

type IroleRepository interface {
	Create(ctx context.Context, role *Role) error
	FindByID(ctx context.Context, id int64) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, size int) ([]Role, int64, error)
	GetByUserID(ctx context.Context, userID int64) ([]Role, error)
	AssignToUser(ctx context.Context, userID, roleID int64) error
	RemoveFromUser(ctx context.Context, userID, roleID int64) error
	GetUserRoleNames(ctx context.Context, userID int64) ([]string, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IroleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) FindByID(ctx context.Context, id int64) (*Role, error) {
	var role Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	return &role, err
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *RoleRepository) Update(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Model(role).Updates(role).Error
}

func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Role{}).Error
}

func (r *RoleRepository) List(ctx context.Context, page, size int) ([]Role, int64, error) {
	var list []Role
	var total int64

	db := r.db.WithContext(ctx).Model(&Role{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := db.Offset(offset).Limit(size).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, total, err
}

func (r *RoleRepository) GetByUserID(ctx context.Context, userID int64) ([]Role, error) {
	var roles []Role
	err := r.db.WithContext(ctx).
		Joins("JOIN usr_user_roles ON usr_user_roles.role_id = usr_roles.id").
		Where("usr_user_roles.user_id = ?", userID).
		Where("usr_user_roles.deleted_at IS NULL").
		Where("usr_roles.status = ?", 1).
		Order("usr_roles.sort_order ASC").
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) AssignToUser(ctx context.Context, userID, roleID int64) error {
	ur := UserRole{UserID: userID, RoleID: roleID}
	return r.db.WithContext(ctx).Create(&ur).Error
}

func (r *RoleRepository) RemoveFromUser(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&UserRole{}).Error
}

func (r *RoleRepository) GetUserRoleNames(ctx context.Context, userID int64) ([]string, error) {
	var names []string
	err := r.db.WithContext(ctx).
		Table("usr_roles").
		Select("usr_roles.name").
		Joins("JOIN usr_user_roles ON usr_user_roles.role_id = usr_roles.id").
		Where("usr_user_roles.user_id = ?", userID).
		Where("usr_user_roles.deleted_at IS NULL").
		Where("usr_roles.status = ?", 1).
		Order("usr_roles.sort_order ASC").
		Pluck("usr_roles.name", &names).Error
	return names, err
}
