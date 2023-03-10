package accrual

import (
	"bonus/internal/model"
	"bonus/internal/service"
)

type AccrualService interface {
	GetOrder(orderNumber string) (Order, error)
}

// делать запросы к accrual по HTTP и сериализовать ответ в структуру accrual.Order и возвращать ее вызывающему коду
type Accrual struct {
	order   service.OrderRepository
	balance service.BalanceRepository
}

func (a *Accrual) Get(order string) {
	//TODO
}

type OrderStatus string

const (
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	//...
)

type Order struct { // или AccrualOrder
	Status OrderStatus
	//...
}

func MapOrderStatus(accrualStatus OrderStatus) model.OrderStatus {
	switch accrualStatus {
	case OrderStatusRegistered, OrderStatusProcessing:
		return model.OrderStatusProcessing
	case OrderStatusProcessed:
		return model.OrderStatusProcessed
	case OrderStatusInvalid:
		return model.OrderStatusInvalid
	}
	return model.OrderStatusInvalid
}
