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
}

func NewUserService(userRepo IuserRepository, infoRepo IuserInfoRepository) *UserService {
	return &UserService{userRepo: userRepo, infoRepo: infoRepo}
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

type UserRoleBrief struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type UserListItem struct {
	*User
	Roles []UserRoleBrief `json:"roles"`
}

type UserListResult struct {
	Total int64          `json:"total"`
	List  []*UserListItem `json:"list"`
}

type UserListReq struct {
	Page    int    `form:"page,default=1" binding:"gte=1"`
	Size    int    `form:"size,default=20" binding:"gte=1,lte=100"`
	Status  *int8  `form:"status"`
	Keyword string `form:"keyword"`
}

func (s *UserService) List(ctx context.Context, req *UserListReq) ([]User, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	return s.userRepo.List(ctx, req.Keyword, req.Status, req.Page, req.Size)
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
