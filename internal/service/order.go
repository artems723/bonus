package service

import (
	"bonus/internal/model"
	"context"
)

type OrderService interface {
}

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByUserID(ctx context.Context, userID uint64) ([]*model.Order, error)
	GetByNumber(ctx context.Context, number string) (*model.Order, error)
}

type orderService struct {
	order OrderRepository
}

func NewOrderService(order OrderRepository) *orderService {
	return &orderService{order}
}
