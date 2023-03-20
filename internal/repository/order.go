package repository

import (
	"bonus/internal/model"
	"bonus/internal/service"
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (o *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	tx := o.db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO orders (user_login, number, status, uploaded_at) VALUES (:user_login, :number, :status, :uploaded_at)", order)
	if err != nil {
		// check if order is already exists
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return service.ErrOrderAlreadyExists
			}
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (o *OrderRepository) GetByLogin(ctx context.Context, login string) ([]*model.Order, error) {
	var orders []*model.Order
	err := o.db.Select(&orders, "SELECT user_login,number,status,accrual,uploaded_at FROM orders WHERE user_login = $1", login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return orders, nil
}
func (o *OrderRepository) GetByNumber(ctx context.Context, number string) (*model.Order, error) {
	var order model.Order
	err := o.db.Get(&order, "SELECT user_login,number,status,accrual,uploaded_at FROM orders WHERE number = $1", number)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return &order, nil
}

func (o *OrderRepository) GetByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	var orders []*model.Order
	err := o.db.Select(&orders, "SELECT user_login,number,status,accrual,uploaded_at FROM orders WHERE status = $1", status)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return orders, nil
}

func (o *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	tx := o.db.MustBegin()
	_, err := tx.NamedExec("UPDATE orders SET status = :status, accrual = :accrual WHERE number = :number", order)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
