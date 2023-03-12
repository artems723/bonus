package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	ID           int64  `json:"id,omitempty" db:"id"`
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"password" db:"password"`
}

func NewUser(login string, passwordHash string) User {
	return User{
		Login:        login,
		PasswordHash: passwordHash,
	}
}

type Order struct {
	ID         int64
	UserID     int64
	Number     string
	Status     OrderStatus
	Accrual    *decimal.Decimal
	UploadedAt time.Time
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Balance struct {
	ID     int64
	UserID int64
	Order  string
	Sum    float64
}
