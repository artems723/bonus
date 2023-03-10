package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	ID           int64
	Login        string
	PasswordHash string
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
