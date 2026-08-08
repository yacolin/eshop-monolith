package staff

import (
	"context"
)

type PermissionService struct {
	permRepo IpermissionRepository
}

func NewPermissionService(permRepo IpermissionRepository) *PermissionService {
	return &PermissionService{permRepo: permRepo}
}

func (s *PermissionService) List(ctx context.Context, category string) ([]*PermissionListItem, error) {
	perms, err := s.permRepo.List(ctx, category)
	if err != nil {
		return nil, err
	}
	items := make([]*PermissionListItem, len(perms))
	for i := range perms {
		items[i] = &PermissionListItem{
			ID:          perms[i].ID,
			Name:        perms[i].Name,
			DisplayName: perms[i].DisplayName,
			Description: perms[i].Description,
			Resource:    perms[i].Resource,
			Action:      perms[i].Action,
			ParentID:    perms[i].ParentID,
			Category:    perms[i].Category,
			SortOrder:   perms[i].SortOrder,
			Status:      perms[i].Status,
		}
	}
	return items, nil
}
