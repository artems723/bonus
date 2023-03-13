package service

import (
	"bonus/internal/model"
	"context"
)

type UserService interface {
	Create(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type userService struct {
	user UserRepository
}

func NewUserService(user UserRepository) *userService {
	return &userService{user}
}

func (u *userService) Create(ctx context.Context, user *model.User) error {
	err := user.HashPassword()
	if err != nil {
		return err
	}
	return u.user.Create(ctx, user)
}
func (u *userService) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	return u.user.GetByLogin(ctx, login)
}
