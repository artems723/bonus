package repository

import (
	"bonus/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type balanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *balanceRepository {
	return &balanceRepository{db}
}

func (b *balanceRepository) Create(ctx context.Context, withdrawal *model.Balance) error {
	return nil
}
func (b *balanceRepository) GetByUserID(ctx context.Context, userID uint64) ([]*model.Balance, error) {
	return nil, nil
}
