package repository

import "github.com/jmoiron/sqlx"

type balanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *balanceRepository {
	return &balanceRepository{db}
}
