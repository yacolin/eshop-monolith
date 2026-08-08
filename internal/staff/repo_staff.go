package staff

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type IstaffRepository interface {
	Create(ctx context.Context, s *Staff) error
	FindByID(ctx context.Context, id int64) (*Staff, error)
	FindByUsername(ctx context.Context, username string) (*Staff, error)
	Update(ctx context.Context, s *Staff) error
	UpdateLastLogin(ctx context.Context, id int64, ip string) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, keyword string, status *int8, page, size int) ([]Staff, int64, error)
}

type StaffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) IstaffRepository {
	return &StaffRepository{db: db}
}

func (r *StaffRepository) Create(ctx context.Context, s *Staff) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *StaffRepository) FindByID(ctx context.Context, id int64) (*Staff, error) {
	var s Staff
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	return &s, err
}

func (r *StaffRepository) FindByUsername(ctx context.Context, username string) (*Staff, error) {
	var s Staff
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&s).Error
	return &s, err
}

func (r *StaffRepository) Update(ctx context.Context, s *Staff) error {
	// 用 map 更新,避免零值字段被跳过(如 status 置 0);密码哈希为空则不更新
	updates := map[string]interface{}{
		"real_name": s.RealName,
		"email":     s.Email,
		"phone":     s.Phone,
		"avatar":    s.Avatar,
		"status":    s.Status,
	}
	if s.PasswordHash != "" {
		updates["password_hash"] = s.PasswordHash
	}
	return r.db.WithContext(ctx).Model(&Staff{}).Where("id = ?", s.ID).Updates(updates).Error
}

func (r *StaffRepository) UpdateLastLogin(ctx context.Context, id int64, ip string) error {
	return r.db.WithContext(ctx).Model(&Staff{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"last_login_at": time.Now(), "last_login_ip": ip}).Error
}

func (r *StaffRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Staff{}).Error
}

func (r *StaffRepository) List(ctx context.Context, keyword string, status *int8, page, size int) ([]Staff, int64, error) {
	db := r.db.WithContext(ctx).Model(&Staff{})
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like, like)
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Staff
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("id ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
