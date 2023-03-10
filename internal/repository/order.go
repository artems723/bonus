package repository

import (
	"bonus/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
	return &orderRepository{db}
}

func (o *orderRepository) Create(ctx context.Context, order *model.Order) error {
	return nil
}
func (o *orderRepository) GetByUserID(ctx context.Context, userID uint64) ([]*model.Order, error) {
	return nil, nil
}
func (o *orderRepository) GetByNumber(ctx context.Context, number string) (*model.Order, error) {
	return nil, nil
}
