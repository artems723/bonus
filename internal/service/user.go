package service

import (
	"bonus/internal/model"
	"context"
)

type UserService interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

type userService struct {
	user UserRepository
}

func NewUserService(user UserRepository) UserService {
	return &userService{
		user: user,
	}
}
