package service

import (
	"bonus/internal/model"
	"context"
)

type AuthService interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

type Auth struct {
}
