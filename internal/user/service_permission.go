package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
)

type PermissionService struct {
	permRepo IpermissionRepository
	roleRepo IroleRepository
}

func NewPermissionService(permRepo IpermissionRepository, roleRepo IroleRepository) *PermissionService {
	return &PermissionService{permRepo: permRepo, roleRepo: roleRepo}
}

func (s *PermissionService) Create(ctx context.Context, req *CreatePermissionReq) (*Permission, error) {
	perm := &Permission{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Category:    req.Category,
	}
	if req.SortOrder != nil {
		perm.SortOrder = *req.SortOrder
	}
	if err := s.permRepo.Create(ctx, perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (s *PermissionService) GetByID(ctx context.Context, id int64) (*Permission, error) {
	perm, err := s.permRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPermissionNotFound
		}
		return nil, err
	}
	return perm, nil
}

func (s *PermissionService) Update(ctx context.Context, id int64, req *UpdatePermissionReq) (*Permission, error) {
	perm, err := s.permRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPermissionNotFound
		}
		return nil, err
	}

	if req.DisplayName != nil {
		perm.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		perm.Description = *req.Description
	}
	if req.Category != nil {
		perm.Category = *req.Category
	}
	if req.SortOrder != nil {
		perm.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		perm.Status = *req.Status
	}

	if err := s.permRepo.Update(ctx, perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (s *PermissionService) Delete(ctx context.Context, id int64) error {
	_, err := s.permRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrPermissionNotFound
		}
		return err
	}
	return s.permRepo.Delete(ctx, id)
}

func (s *PermissionService) List(ctx context.Context, req *PermissionListReq) (*PermissionListResult, error) {
	req.Normalize()
	list, total, err := s.permRepo.List(ctx, req.Category, req.Resource, req.RoleID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*Permission, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &PermissionListResult{Total: total, List: items}, nil
}

func (s *PermissionService) CheckUserPermissions(ctx context.Context, userID int64, permissionNames []string) (map[string]bool, error) {
	roles, err := s.roleRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]int64, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.ID
	}

	result := make(map[string]bool, len(permissionNames))
	for _, name := range permissionNames {
		has, err := s.permRepo.HasPermissionByRoleIDs(ctx, roleIDs, name)
		if err != nil {
			return nil, err
		}
		result[name] = has
	}
	return result, nil
}

func (s *PermissionService) AssignToRole(ctx context.Context, roleID int64, permissionIDs []int64) error {
	for _, permID := range permissionIDs {
		if err := s.permRepo.AssignToRole(ctx, roleID, permID); err != nil {
			return err
		}
	}
	return nil
}

func (s *PermissionService) RemoveFromRole(ctx context.Context, roleID int64, permissionIDs []int64) error {
	for _, permID := range permissionIDs {
		if err := s.permRepo.RemoveFromRole(ctx, roleID, permID); err != nil {
			return err
		}
	}
	return nil
}

func (s *PermissionService) GetByRoleID(ctx context.Context, roleID int64) ([]Permission, error) {
	return s.permRepo.GetByRoleID(ctx, roleID)
}
