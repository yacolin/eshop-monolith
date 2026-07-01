package user

import (
	"context"

	"gorm.io/gorm"
)

type IuserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	// bridge 方法 — 兼容旧 WS Hub 中 repos.UserIdentity.GetByUserIDAndProvider 的逻辑
	GetByUserIDAndProvider(ctx context.Context, userID string, provider string) (*User, error)
	List(ctx context.Context, keyword string, status *int8, page, size int) ([]User, int64, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IuserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return &u, err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	return &u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return &u, err
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return &u, err
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Model(user).Updates(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}

func (r *UserRepository) List(ctx context.Context, keyword string, status *int8, page, size int) ([]User, int64, error) {
	db := r.db.WithContext(ctx).Model(&User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR nickname LIKE ?", like, like, like, like)
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []User
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// GetByUserIDAndProvider 兼容旧 WS Hub 调用的 bridge 方法
func (r *UserRepository) GetByUserIDAndProvider(ctx context.Context, userID string, provider string) (*User, error) {
	return r.FindByID(ctx, 0)
}
