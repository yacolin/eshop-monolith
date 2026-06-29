package user

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

func (s *RoleService) Create(ctx context.Context, req *CreateRoleReq) (*Role, error) {
	role := &Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsSystem:    req.IsSystem,
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if req.SortOrder != nil {
		role.SortOrder = *req.SortOrder
	}
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetByID(ctx context.Context, id int64) (*Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetByName(ctx context.Context, name string) (*Role, error) {
	role, err := s.roleRepo.FindByName(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Update(ctx context.Context, id int64, req *UpdateRoleReq) (*Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}

	if role.IsSystem {
		return nil, errcode.ErrCannotModifySystemRole
	}

	if req.DisplayName != nil {
		role.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if req.SortOrder != nil {
		role.SortOrder = *req.SortOrder
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Delete(ctx context.Context, id int64) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	if role.IsSystem {
		return errcode.ErrCannotDeleteSystemRole
	}
	return s.roleRepo.Delete(ctx, id)
}

func (s *RoleService) List(ctx context.Context, page, size int) (*RoleListResult, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	list, total, err := s.roleRepo.List(ctx, page, size)
	if err != nil {
		return nil, err
	}
	items := make([]*Role, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &RoleListResult{Total: total, List: items}, nil
}

func (s *RoleService) AssignToUser(ctx context.Context, userID, roleID int64) error {
	return s.roleRepo.AssignToUser(ctx, userID, roleID)
}

func (s *RoleService) RemoveFromUser(ctx context.Context, userID, roleID int64) error {
	return s.roleRepo.RemoveFromUser(ctx, userID, roleID)
}

func (s *RoleService) GetUserRoles(ctx context.Context, userID int64) ([]Role, error) {
	return s.roleRepo.GetByUserID(ctx, userID)
}
