package repository

import (
	"bonus/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type BalanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{db}
}

func (b *BalanceRepository) Create(ctx context.Context, withdrawal *model.Balance) error {
	return nil
}
func (b *BalanceRepository) GetByUserID(ctx context.Context, userID uint64) ([]*model.Balance, error) {
	return nil, nil
}
