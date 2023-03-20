package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Balance struct {
	UserLogin   string `db:"user_login"`
	OrderNumber string `db:"order_number"`
	Debit       *decimal.Decimal
	Credit      *decimal.Decimal
	ProcessedAt time.Time `db:"processed_at"`
}

type CurrentBalance struct {
	UserLogin string `json:"-"`
	Current   *decimal.Decimal
	Withdrawn *decimal.Decimal
}

type Withdrawal struct {
	Order       string           `json:"order" db:"order_number"`
	Sum         *decimal.Decimal `json:"sum" db:"credit"`
	ProcessedAt time.Time        `json:"processed_at,omitempty" db:"processed_at"`
}
