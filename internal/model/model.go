package model

import "time"

type User struct {
	ID           uint64
	Login        string
	PasswordHash string
}

type Order struct {
	ID         uint64
	UserID     uint64
	Number     string
	Status     string
	Accrual    float64
	UploadedAt time.Time
}

type Balance struct {
	ID     uint64
	UserID uint64
	Order  string
	Sum    float64
}
