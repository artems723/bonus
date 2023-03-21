package service

import (
	"bonus/internal/model"
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

type BalanceService interface {
	GetByLogin(ctx context.Context, login string) (*model.CurrentBalance, error)
	Withdraw(ctx context.Context, login string, withdrawal *model.Withdrawal) error
	GetWithdrawals(ctx context.Context, login string) ([]*model.Withdrawal, error)
}

type BalanceRepository interface {
	Create(ctx context.Context, balance *model.Balance) error
	GetByLogin(ctx context.Context, login string) ([]*model.Balance, error)
	GetWithdrawals(ctx context.Context, login string) ([]*model.Withdrawal, error)
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
		if bal.Debit != nil {
			current = current.Add(*bal.Debit)
		}
		if bal.Credit != nil {
			current = current.Sub(*bal.Credit)
			withdrawn = withdrawn.Add(*bal.Credit)
		}
	}

	return &currentBalance, nil
}

func (b *balanceService) Withdraw(ctx context.Context, login string, withdrawal *model.Withdrawal) error {
	currentBalance, err := b.GetByLogin(ctx, login)
	if err != nil {
		return err
	}
	if currentBalance.Current.Cmp(*withdrawal.Sum) < 0 {
		return ErrNotEnoughFunds
	}

	balance := model.Balance{
		UserLogin:   login,
		OrderNumber: withdrawal.Order,
		Debit:       nil,
		Credit:      withdrawal.Sum,
		ProcessedAt: time.Now(),
	}

	err = b.balance.Create(ctx, &balance)
	if err != nil {
		return err
	}
	return nil
}

func (b *balanceService) GetWithdrawals(ctx context.Context, login string) ([]*model.Withdrawal, error) {
	withdrawals, err := b.balance.GetWithdrawals(ctx, login)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
