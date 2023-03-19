package service

import (
	"bonus/internal/model"
	"context"
	"errors"
	"github.com/shopspring/decimal"
)

type BalanceService interface {
	GetByLogin(ctx context.Context, login string) (*model.CurrentBalance, error)
}

type BalanceRepository interface {
	Create(ctx context.Context, withdrawal *model.Balance) error
	GetByLogin(ctx context.Context, login string) ([]*model.Balance, error)
}

type balanceService struct {
	balance BalanceRepository
}

func NewBalanceService(balance BalanceRepository) *balanceService {
	return &balanceService{balance}
}

func (b *balanceService) GetByLogin(ctx context.Context, login string) (*model.CurrentBalance, error) {
	current := decimal.NewFromInt(0)
	withdrawn := decimal.NewFromInt(0)
	currentBalance := model.CurrentBalance{
		UserLogin: login,
		Current:   &current,
		Withdrawn: &withdrawn,
	}

	balances, err := b.balance.GetByLogin(ctx, login)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if errors.Is(err, ErrNotFound) {
		return &currentBalance, nil
	}

	for _, bal := range balances {
		current = current.Add(*bal.Debit)
		current = current.Sub(*bal.Credit)
		withdrawn = withdrawn.Add(*bal.Credit)
	}

	return &currentBalance, nil
}
