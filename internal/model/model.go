package model

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"password" db:"password_hash"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		Password string `json:"password"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	u.PasswordHash, err = HashPassword(aux.Password)
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
