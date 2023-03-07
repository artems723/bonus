package service

type AccrualService interface {
	Get(order string)
}

type Accrual struct {
	order   OrderRepository
	balance BalanceRepository
}

func (a *Accrual) Get(order string) {
	//TODO
}
