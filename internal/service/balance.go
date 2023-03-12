package service

import (
	"bonus/internal/model"
	"context"
)

type BalanceService interface {
}

type BalanceRepository interface {
	Create(ctx context.Context, withdrawal *model.Balance) error
	GetByUserID(ctx context.Context, userID uint64) ([]*model.Balance, error)
}

type balanceService struct {
	balance BalanceRepository
}

func NewBalanceService(balance BalanceRepository) *balanceService {
	return &balanceService{balance}
}
