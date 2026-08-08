package staff

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
)

type RoleService struct {
	roleRepo IroleRepository
}

func NewRoleService(roleRepo IroleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) List(ctx context.Context) ([]*RoleListItem, error) {
	roles, err := s.roleRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*RoleListItem, len(roles))
	for i := range roles {
		ids, err := s.roleRepo.GetPermissionIDsByRoleID(ctx, roles[i].ID)
		if err != nil {
			return nil, err
		}
		items[i] = &RoleListItem{
			ID:              roles[i].ID,
			Name:            roles[i].Name,
			DisplayName:     roles[i].DisplayName,
			Description:     roles[i].Description,
			RoleType:        roles[i].RoleType,
			Status:          roles[i].Status,
			SortOrder:       roles[i].SortOrder,
			PermissionCount: int64(len(ids)),
		}
	}
	return items, nil
}

func (s *RoleService) Create(ctx context.Context, req *CreateRoleReq) (*SysRole, error) {
	role := &SysRole{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		RoleType:    "custom",
		SortOrder:   req.SortOrder,
		Status:      1,
	}
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

// Update 更新自定义角色;builtin 角色不可修改(ErrCannotModifySystemRole)
func (s *RoleService) Update(ctx context.Context, roleID int64, req *UpdateRoleReq) error {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	if role.RoleType == "builtin" {
		return errcode.ErrCannotModifySystemRole
	}
	role.DisplayName = req.DisplayName
	role.Description = req.Description
	role.SortOrder = req.SortOrder
	if req.Status != nil {
		role.Status = *req.Status
	}
	return s.roleRepo.Update(ctx, role)
}

// Delete 删除自定义角色;builtin 角色不可删(ErrCannotDeleteSystemRole)
func (s *RoleService) Delete(ctx context.Context, roleID int64) error {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	if role.RoleType == "builtin" {
		return errcode.ErrCannotDeleteSystemRole
	}
	return s.roleRepo.Delete(ctx, roleID)
}

func (s *RoleService) GetPermissionIDs(ctx context.Context, roleID int64) ([]int64, error) {
	return s.roleRepo.GetPermissionIDsByRoleID(ctx, roleID)
}

func (s *RoleService) AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	return s.roleRepo.ReplaceRolePermissions(ctx, roleID, permissionIDs)
}
