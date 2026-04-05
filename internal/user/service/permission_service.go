package service

import (
	"context"
	"errors"
	"eshop-monolith/internal/pkg/errcode"
	"eshop-monolith/internal/user/api/dto"
	"eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/user/domain/repositories"

	"gorm.io/gorm"
)

type PermissionService interface {
	CreatePermission(req *dto.CreatePermissionRequest) (*models.Permission, error)
	GetPermission(id int64) (*models.Permission, error)
	GetPermissionByName(name string) (*models.Permission, error)
	UpdatePermission(id int64, req *dto.UpdatePermissionRequest) (*models.Permission, error)
	DeletePermission(id int64) error

	ListPermissions(page, pageSize int) (*dto.ListPermissionsResponse, error)
	GetPermissionsByCategory(category string, page, pageSize int) (*dto.ListPermissionsResponse, error)
	GetPermissionsByResource(resource string, page, pageSize int) (*dto.ListPermissionsResponse, error)
	GetPermissionsByRoleID(roleID int64, page, pageSize int) (*dto.ListPermissionsResponse, error)

	CheckPermissionsByRoleIDs(roleIDs []int64, permissionNames []string) (map[string]bool, error)
	CheckUserPermissions(userID int64, permissionNames []string) (map[string]bool, error)

	CreateRole(req *dto.CreateRoleRequest) (*models.Role, error)
	GetRole(id int64) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	UpdateRole(id int64, req *dto.UpdateRoleRequest) (*models.Role, error)
	DeleteRole(id int64) error
	ListRoles(page, pageSize int) (*dto.ListRolesResponse, error)
	AssignRoleToUser(userID int64, roleID int64) error
	RemoveRoleFromUser(userID int64, roleID int64) error
	GetUserRoles(userID int64) ([]models.Role, error)
	AssignPermissionsToRole(roleID int64, permissionIDs []int64) error
	RemovePermissionsFromRole(roleID int64, permissionIDs []int64) error
}

type permissionService struct {
	permissionRepo repositories.IpermissionRepository
	userRepo       repositories.IuserRepository
	roleRepo       repositories.IroleRepository
}

func NewPermissionService(
	permissionRepo repositories.IpermissionRepository,
	userRepo repositories.IuserRepository,
	roleRepo repositories.IroleRepository,
) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
	}
}

func (s *permissionService) CreatePermission(req *dto.CreatePermissionRequest) (*models.Permission, error) {
	// 检查权限名称是否已存在
	exists, err := s.permissionRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("permission name already exists")
	}

	permission := &models.Permission{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Category:    req.Category,
		Sort:        req.Sort,
		Status:      1,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) GetPermission(id int64) (*models.Permission, error) {
	return s.permissionRepo.GetByID(id)
}

func (s *permissionService) GetPermissionByName(name string) (*models.Permission, error) {
	return s.permissionRepo.GetByName(name)
}

func (s *permissionService) UpdatePermission(id int64, req *dto.UpdatePermissionRequest) (*models.Permission, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	if req.DisplayName != nil {
		permission.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		permission.Description = *req.Description
	}
	if req.Category != nil {
		permission.Category = *req.Category
	}
	if req.Sort != nil {
		permission.Sort = *req.Sort
	}
	if req.Status != nil {
		permission.Status = *req.Status
	}

	if err := s.permissionRepo.Update(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) DeletePermission(id int64) error {
	// TODO: 检查是否有角色使用此权限
	return s.permissionRepo.Delete(id)
}

func (s *permissionService) ListPermissions(page, pageSize int) (*dto.ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.List(pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &dto.ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func (s *permissionService) GetPermissionsByCategory(category string, page, pageSize int) (*dto.ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.ByCategory(category, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &dto.ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

// Add dto.ListPermissionsResponse
func (s *permissionService) GetPermissionsByResource(resource string, page, pageSize int) (*dto.ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.ByResource(resource, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// 添加dto前缀
	return &dto.ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

// Add dto.ListPermissionsResponse
func (s *permissionService) GetPermissionsByRoleID(roleID int64, page, pageSize int) (*dto.ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	permissions, total, err := s.permissionRepo.ByRoleID(roleID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// 添加dto前缀
	return &dto.ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func (s *permissionService) CheckPermissionsByRoleIDs(roleIDs []int64, permissionNames []string) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, permissionName := range permissionNames {
		has, err := s.permissionRepo.HasPermissionByRoleIDs(roleIDs, permissionName)
		if err != nil {
			return nil, err
		}
		result[permissionName] = has
	}

	return result, nil
}

func (s *permissionService) CheckUserPermissions(userID int64, permissionNames []string) (map[string]bool, error) {
	roles, err := s.roleRepo.GetUserRoles(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]int64, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	return s.CheckPermissionsByRoleIDs(roleIDs, permissionNames)
}

func (s *permissionService) CreateRole(req *dto.CreateRoleRequest) (*models.Role, error) {
	role := &models.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Status:      req.Status,
		Sort:        req.Sort,
		IsSystem:    req.IsSystem,
	}

	if role.Status == 0 {
		role.Status = 1
	}

	if err := s.roleRepo.Create(context.Background(), role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *permissionService) GetRole(id int64) (*models.Role, error) {
	return s.roleRepo.GetByID(context.Background(), id)
}

func (s *permissionService) GetRoleByName(name string) (*models.Role, error) {
	return s.roleRepo.GetByName(context.Background(), name)
}

func (s *permissionService) UpdateRole(id int64, req *dto.UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(context.Background(), id)
	if err != nil {
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
	if req.Sort != nil {
		role.Sort = *req.Sort
	}

	if err := s.roleRepo.Update(context.Background(), role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *permissionService) DeleteRole(id int64) error {
	role, err := s.roleRepo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return errcode.ErrCannotDeleteSystemRole
	}

	return s.roleRepo.Delete(context.Background(), id)
}

func (s *permissionService) ListRoles(page, pageSize int) (*dto.ListRolesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	roles, total, err := s.roleRepo.List(context.Background(), pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &dto.ListRolesResponse{
		Roles:    roles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *permissionService) AssignRoleToUser(userID int64, roleID int64) error {
	return s.roleRepo.AssignRoleToUser(context.Background(), userID, roleID)
}

func (s *permissionService) RemoveRoleFromUser(userID int64, roleID int64) error {
	return s.roleRepo.RemoveRoleFromUser(context.Background(), userID, roleID)
}

func (s *permissionService) GetUserRoles(userID int64) ([]models.Role, error) {
	return s.roleRepo.GetUserRoles(context.Background(), userID)
}

func (s *permissionService) AssignPermissionsToRole(roleID int64, permissionIDs []int64) error {
	for _, permissionID := range permissionIDs {
		if err := s.roleRepo.AssignPermissionToRole(context.Background(), roleID, permissionID); err != nil {
			return err
		}
	}
	return nil
}

func (s *permissionService) RemovePermissionsFromRole(roleID int64, permissionIDs []int64) error {
	for _, permissionID := range permissionIDs {
		if err := s.roleRepo.RemovePermissionFromRole(context.Background(), roleID, permissionID); err != nil {
			return err
		}
	}
	return nil
}
