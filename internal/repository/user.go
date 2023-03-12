package repository

import (
	"bonus/internal/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) Create(ctx context.Context, user *model.User) error {
	tx := u.db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO users (login, password_hash) VALUES (:login, :password)", user)
	if err != nil {
		// check if login is taken
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrLoginIsTaken
			}
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	return &model.User{}, nil
}

var ErrUserNotFound = errors.New("user not found")
var ErrLoginIsTaken = errors.New("username is taken, try another one")
