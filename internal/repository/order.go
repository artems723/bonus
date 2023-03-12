package repository

import (
	"bonus/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (o *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	return nil
}
func (o *OrderRepository) GetByUserID(ctx context.Context, userID uint64) ([]*model.Order, error) {
	return nil, nil
}
func (o *OrderRepository) GetByNumber(ctx context.Context, number string) (*model.Order, error) {
	return nil, nil
}
