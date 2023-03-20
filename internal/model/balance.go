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
	CreatedAt   time.Time `db:"created_at"`
}

type CurrentBalance struct {
	UserLogin string `json:"-"`
	Current   *decimal.Decimal
	Withdrawn *decimal.Decimal
}

type Withdrawal struct {
	Order string
	Sum   *decimal.Decimal
}
