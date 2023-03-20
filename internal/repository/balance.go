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
	_, err := tx.NamedExec("INSERT INTO balances (user_login, order_number, debit, credit, processed_at) VALUES (:user_login, :order_number, :debit, :credit, :processed_at)", balance)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return nil
}
func (b *BalanceRepository) GetByLogin(ctx context.Context, login string) ([]*model.Balance, error) {
	var balances []*model.Balance
	err := b.db.Select(&balances, "SELECT user_login,order_number,debit,credit,processed_at FROM balances WHERE user_login = $1", login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return balances, nil
}

func (b *BalanceRepository) GetWithdrawals(ctx context.Context, login string) ([]*model.Withdrawal, error) {
	var withdrawals []*model.Withdrawal
	err := b.db.Select(&withdrawals, "SELECT order_number,credit,processed_at FROM balances WHERE user_login = $1 AND credit IS NOT NULL", login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return withdrawals, nil
}
