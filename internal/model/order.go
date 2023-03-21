package model

import (
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type Order struct {
	UserLogin  string           `json:"-" db:"user_login"`
	Number     string           `json:"number" db:"number"`
	Status     OrderStatus      `json:"status" db:"status"`
	Accrual    *decimal.Decimal `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time        `json:"uploaded_at" db:"uploaded_at"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

func (o Order) Valid() bool {
	number, err := strconv.Atoi(o.Number)
	if err != nil {
		return false
	}
	return checkLuhn(number)
}

// Check number for Luhn algorithm
func checkLuhn(number int) bool {
	number = number / 10
	var luhn int
	for i := 0; number > 0; i++ {
		cur := number % 10
		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}
		luhn += cur
		number = number / 10
	}
	return (number%10+luhn%10)%10 == 0
}
