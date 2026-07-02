package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
)

type UserService struct {
	userRepo IuserRepository
	infoRepo IuserInfoRepository
	roleRepo IroleRepository
}

func NewUserService(userRepo IuserRepository, infoRepo IuserInfoRepository, roleRepo IroleRepository) *UserService {
	return &UserService{userRepo: userRepo, infoRepo: infoRepo, roleRepo: roleRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*User, *UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errcode.ErrUserNotFound
		}
		return nil, nil, err
	}

	info, err := s.infoRepo.GetUserInfoByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	return user, info, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	info, err := s.infoRepo.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return info, nil
}

func (s *UserService) List(ctx context.Context, req *UserListReq) (*UserListResult, error) {
	users, total, err := s.userRepo.List(ctx, req.Keyword, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*UserListItem, len(users))
	for i := range users {
		roles, _ := s.roleRepo.GetByUserID(ctx, users[i].ID)
		roleBriefs := make([]UserRoleBrief, len(roles))
		for j := range roles {
			roleBriefs[j] = UserRoleBrief{ID: roles[j].ID, Name: roles[j].Name, DisplayName: roles[j].DisplayName}
		}
		items[i] = &UserListItem{User: &users[i], Roles: roleBriefs}
	}
	return &UserListResult{Total: total, List: items}, nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, userID int64, req *UpdateUserInfoReq) (*UserInfo, error) {
	info, err := s.infoRepo.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			info = &UserInfo{UserID: userID}
		} else {
			return nil, err
		}
	}

	if req.Gender != 0 {
		info.Gender = req.Gender
	}
	if req.Bio != "" {
		info.Bio = req.Bio
	}
	if req.Country != "" {
		info.Country = req.Country
	}
	if req.Province != "" {
		info.Province = req.Province
	}
	if req.City != "" {
		info.City = req.City
	}
	if req.ZipCode != "" {
		info.ZipCode = req.ZipCode
	}
	if req.Language != "" {
		info.Language = req.Language
	}
	if req.Timezone != "" {
		info.Timezone = req.Timezone
	}

	if info.ID == 0 {
		if err := s.infoRepo.CreateUserInfo(ctx, info); err != nil {
			return nil, err
		}
	} else {
		if err := s.infoRepo.UpdateUserInfo(ctx, info); err != nil {
			return nil, err
		}
	}

	return info, nil
}
