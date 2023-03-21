package repository

import (
	"bonus/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrLoginIsTaken = errors.New("username is taken, try another one")
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) Create(ctx context.Context, user *model.User) error {
	tx := u.db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO users (login, password_hash) VALUES (:login, :password_hash)", user)
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
	var user model.User
	err := u.db.Get(&user, "SELECT login, password_hash FROM users WHERE login = $1", login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return &user, nil
}
