package service

type AccrualService interface {
}

type Accrual struct {
	order   OrderRepository
	balance BalanceRepository
}