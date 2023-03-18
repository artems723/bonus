package service

import (
	"bonus/internal/model"
	"context"
	"errors"
)

type OrderService interface {
	Create(ctx context.Context, order *model.Order) error
	GetByLogin(ctx context.Context, login string) ([]*model.Order, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByLogin(ctx context.Context, login string) ([]*model.Order, error)
	GetByNumber(ctx context.Context, number string) (*model.Order, error)
}

type orderService struct {
	order OrderRepository
}

func NewOrderService(order OrderRepository) *orderService {
	return &orderService{
		order: order,
	}
}

func (s *orderService) Create(ctx context.Context, order *model.Order) error {
	err := s.order.Create(ctx, order)
	if err != nil && !errors.Is(err, ErrOrderAlreadyExists) {
		return err
	}
	if errors.Is(err, ErrOrderAlreadyExists) {
		currOrder, err2 := s.order.GetByNumber(ctx, order.Number)
		if err2 != nil {
			return err2
		}
		if currOrder.UserLogin != order.UserLogin {
			return ErrOrderAlreadyExistsForAnotherUser
		}
		return err
	}
	return nil
}

func (s *orderService) GetByLogin(ctx context.Context, login string) ([]*model.Order, error) {
	orders, err := s.order.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

var ErrOrderAlreadyExistsForAnotherUser = errors.New("order already exists for another user")
var ErrOrderAlreadyExists = errors.New("order already exists")
var ErrNotFound = errors.New("not found")
