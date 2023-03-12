package repository

import (
	"bonus/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db}
}

func (u *userRepository) Create(ctx context.Context, user model.User) error {
	return nil
}
func (u *userRepository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	return model.User{}, nil
}
