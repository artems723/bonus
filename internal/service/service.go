package service

type Service interface {
}

type service struct {
	order   OrderRepository
	balance BalanceRepository
}

func New(order OrderRepository, balance BalanceRepository) Service {
	return &service{
		order:   order,
		balance: balance,
	}
}
