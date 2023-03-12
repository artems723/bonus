package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"password" db:"password"`
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
