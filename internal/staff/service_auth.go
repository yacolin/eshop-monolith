package staff

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/token"
	"eshop-monolith/pkg/utils"
)

type AuthService struct {
	staffRepo        IstaffRepository
	staffRoleRepo    IstaffRoleRepository
	loginHistoryRepo IstaffLoginHistoryRepository
}

func NewAuthService(staffRepo IstaffRepository, staffRoleRepo IstaffRoleRepository, loginHistoryRepo IstaffLoginHistoryRepository) *AuthService {
	return &AuthService{
		staffRepo:        staffRepo,
		staffRoleRepo:    staffRoleRepo,
		loginHistoryRepo: loginHistoryRepo,
	}
}

// Login 员工密码登录:校验凭据 → 拉角色 → 签发令牌 → 记录登录历史(成功/失败都写) → 更新最后登录
func (s *AuthService) Login(ctx context.Context, username, password, ip, userAgent string) (*Staff, []string, *token.TokenPair, error) {
	if username == "" || password == "" {
		return nil, nil, nil, errcode.ErrInvalidCredentials
	}

	staff, err := s.staffRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.recordHistory(ctx, 0, ip, userAgent, "password", 0, "invalid credentials")
			return nil, nil, nil, errcode.ErrInvalidCredentials
		}
		return nil, nil, nil, err
	}

	if staff.Status != 1 {
		s.recordHistory(ctx, staff.ID, ip, userAgent, "password", 0, "account disabled")
		return nil, nil, nil, errcode.ErrAccountDisabled
	}

	if !utils.CheckPasswordHash(password, staff.PasswordHash) {
		s.recordHistory(ctx, staff.ID, ip, userAgent, "password", 0, "wrong password")
		return nil, nil, nil, errcode.ErrInvalidCredentials
	}

	roleNames, err := s.staffRoleRepo.GetRoleNamesByStaffID(ctx, staff.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	tokenPair, err := token.GenerateTokenPair(staff.ID, staff.Username, roleNames)
	if err != nil {
		return nil, nil, nil, err
	}

	s.recordHistory(ctx, staff.ID, ip, userAgent, "password", 1, "")
	_ = s.staffRepo.UpdateLastLogin(ctx, staff.ID, ip)
	return staff, roleNames, tokenPair, nil
}

// Profile 查询员工信息与角色
func (s *AuthService) Profile(ctx context.Context, staffID int64) (*Staff, []string, error) {
	staff, err := s.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errcode.ErrNotFound
		}
		return nil, nil, err
	}
	roleNames, err := s.staffRoleRepo.GetRoleNamesByStaffID(ctx, staffID)
	if err != nil {
		return nil, nil, err
	}
	return staff, roleNames, nil
}

func (s *AuthService) recordHistory(ctx context.Context, staffID int64, ip, device, method string, status int8, reason string) {
	go func() {
		_ = s.loginHistoryRepo.Create(context.Background(), &StaffLoginHistory{
			StaffID:       staffID,
			LoginIP:       ip,
			LoginDevice:   device,
			LoginMethod:   method,
			LoginStatus:   status,
			FailureReason: reason,
		})
	}()
}
