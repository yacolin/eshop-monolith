package staff

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/utils"
)

type StaffService struct {
	staffRepo     IstaffRepository
	staffRoleRepo IstaffRoleRepository
}

func NewStaffService(staffRepo IstaffRepository, staffRoleRepo IstaffRoleRepository) *StaffService {
	return &StaffService{staffRepo: staffRepo, staffRoleRepo: staffRoleRepo}
}

func (s *StaffService) List(ctx context.Context, req *StaffListReq) (*StaffListResult, error) {
	req.Normalize()
	staffs, total, err := s.staffRepo.List(ctx, req.Keyword, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*StaffListItem, len(staffs))
	for i := range staffs {
		roles, _ := s.staffRoleRepo.GetRoleNamesByStaffID(ctx, staffs[i].ID)
		items[i] = &StaffListItem{
			ID:          staffs[i].ID,
			Username:    staffs[i].Username,
			RealName:    staffs[i].RealName,
			Email:       staffs[i].Email,
			Phone:       staffs[i].Phone,
			Status:      staffs[i].Status,
			LastLoginAt: staffs[i].LastLoginAt,
			Roles:       roles,
		}
	}
	return &StaffListResult{Total: total, List: items}, nil
}

func (s *StaffService) Create(ctx context.Context, req *CreateStaffReq, operatorID int64) (*Staff, error) {
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	staff := &Staff{
		Username:     req.Username,
		PasswordHash: string(hash),
		RealName:     req.RealName,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       status,
		CreatedBy:    operatorID,
	}
	if err := s.staffRepo.Create(ctx, staff); err != nil {
		return nil, err
	}
	return staff, nil
}

// Update 更新员工资料/状态/密码。保护规则:不可操作自己;内置 admin(id=1)不可禁用。
func (s *StaffService) Update(ctx context.Context, staffID int64, req *UpdateStaffReq, operatorID int64) error {
	if staffID == operatorID {
		return errcode.ErrCannotOperateSelf
	}
	staff, err := s.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	staff.RealName = req.RealName
	staff.Email = req.Email
	staff.Phone = req.Phone
	staff.Avatar = req.Avatar
	if req.Status != nil {
		if staff.ID == 1 && *req.Status != 1 {
			return errcode.ErrCannotOperateSelf
		}
		staff.Status = *req.Status
	}
	if req.Password != "" {
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		staff.PasswordHash = string(hash)
	}
	return s.staffRepo.Update(ctx, staff)
}

// Delete 删除员工(软删)。保护规则:不可操作自己;内置 admin(id=1)不可删。
func (s *StaffService) Delete(ctx context.Context, staffID, operatorID int64) error {
	if staffID == operatorID || staffID == 1 {
		return errcode.ErrCannotOperateSelf
	}
	return s.staffRepo.Delete(ctx, staffID)
}

func (s *StaffService) AssignRoles(ctx context.Context, staffID int64, roleIDs []int64) error {
	return s.staffRoleRepo.ReplaceStaffRoles(ctx, staffID, roleIDs)
}
