package service

import (
	"context"
	"fmt"
	"strconv"

	"eshop-monolith/internal/domain/user"
	"eshop-monolith/internal/eventbus"
)

// UserService 用户服务
type UserService struct {
	repo user.Repository
	bus  *eventbus.Bus
}

// NewUserService 创建用户服务
func NewUserService(repo user.Repository, bus *eventbus.Bus) *UserService {
	return &UserService{
		repo: repo,
		bus:  bus,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Email string `json:"email"`
}

// Register 注册用户
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*user.User, error) {
	// 创建用户
	newUser := &user.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 密码应该在调用前已经加密
		Role:     "user",
	}

	// 保存用户
	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// 发布用户注册事件
	s.bus.Publish(user.UserRegisteredEvent{
		UserID:   fmt.Sprintf("%d", newUser.ID),
		Username: newUser.Username,
		Email:    newUser.Email,
	})

	return newUser, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest, ip string) (*user.User, error) {
	// 根据用户名查找用户
	foundUser, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// 发布用户登录事件
	s.bus.Publish(user.UserLoggedInEvent{
		UserID:   fmt.Sprintf("%d", foundUser.ID),
		Username: foundUser.Username,
		IP:       ip,
	})

	return foundUser, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, idInt)
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(ctx context.Context, id string, req *UpdateProfileRequest) (*user.User, error) {
	// 获取用户
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	existingUser, err := s.repo.FindByID(ctx, idInt)
	if err != nil {
		return nil, err
	}

	// 更新资料
	existingUser.Email = req.Email

	// 保存用户
	if err := s.repo.Update(ctx, existingUser); err != nil {
		return nil, err
	}

	// 发布用户资料更新事件
	s.bus.Publish(user.UserProfileUpdatedEvent{
		UserID:   fmt.Sprintf("%d", existingUser.ID),
		Username: existingUser.Username,
	})

	return existingUser, nil
}

// ListRoles 列出所有角色
func (s *UserService) ListRoles(ctx context.Context) ([]user.Role, error) {
	return s.repo.ListRoles(ctx)
}

// ListPermissions 列出所有权限
func (s *UserService) ListPermissions(ctx context.Context) ([]user.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

// CheckPermission 检查用户权限
func (s *UserService) CheckPermission(ctx context.Context, userID string, permissionName string) (bool, error) {
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return false, err
	}
	return s.repo.CheckPermission(ctx, userIDInt, permissionName)
}
