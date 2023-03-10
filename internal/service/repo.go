package service

import (
	"bonus/internal/model"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByUserID(ctx context.Context, userID uint64) ([]*model.Order, error)
	GetByNumber(ctx context.Context, number string) (*model.Order, error)
}

type BalanceRepository interface {
	Create(ctx context.Context, withdrawal *model.Balance) error
	GetByUserID(ctx context.Context, userID uint64) ([]*model.Balance, error)
}
