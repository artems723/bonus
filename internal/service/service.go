package service

type Service interface {
	Shutdown() error
}

type OrderRepository interface {
}

type BalanceRepository interface {
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

func (s *service) Shutdown() error {
	return nil
}
