package repository

import (
	"bonus/internal/model"
	"bonus/internal/service"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type BalanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{db}
}

func (b *BalanceRepository) Create(ctx context.Context, balance *model.Balance) error {
	tx := b.db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO balances (user_login, order_number, debit, credit, created_at) VALUES (:user_login, :order_number, :debit, :credit, :created_at)", balance)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return nil
}
func (b *BalanceRepository) GetByLogin(ctx context.Context, login string) ([]*model.Balance, error) {
	var balances []*model.Balance
	err := b.db.Select(&balances, "SELECT user_login,order_number,debit,credit,created_at FROM balances WHERE user_login = $1", login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return balances, nil
}
