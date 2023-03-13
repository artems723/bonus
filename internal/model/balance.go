package model

type Balance struct {
	ID     int64
	UserID int64
	Order  string
	Sum    float64
}
