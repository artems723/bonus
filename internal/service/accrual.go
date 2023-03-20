package service

import (
	"bonus/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"log"
	"net/http"
	"time"
)

type AccrualService interface {
	GetOrder(orderNumber string) (*AccrualOrder, error)
}

// делать запросы к accrual по HTTP и сериализовать ответ в структуру accrual.Order и возвращать ее вызывающему коду
type Accrual struct {
	order   OrderRepository
	balance BalanceRepository
	address string
}

func NewAccrualService(order OrderRepository, balance BalanceRepository, address string) *Accrual {
	return &Accrual{
		order:   order,
		balance: balance,
		address: address,
	}
}

func (a *Accrual) Run(ctx context.Context) {
	a.HandleStatusNew(ctx)
	a.HandleStatusProcessing(ctx)
}

func (a *Accrual) HandleStatusNew(ctx context.Context) {
	// Get all orders with status NEW
	orders, err := a.order.GetByStatus(ctx, model.OrderStatusNew)
	if err != nil {
		log.Printf("error getting orders with status NEW: %v", err)
		return
	}
	log.Printf("orders with status NEW: %v", orders)
	// Get info from accrual service for all orders
	for _, order := range orders {
		accrualOrder, err := a.GetOrder(order.Number)
		if err != nil {
			log.Printf("error getting order from accrual service: %v", err)
			continue
		}
		log.Printf("order from accrual service: %v", accrualOrder)
		// Update order status in DB
		order.Status = MapOrderStatus(accrualOrder.Status)
		if err := a.order.Update(ctx, order); err != nil {
			log.Printf("error updating order status: %v", err)
			continue
		}
		log.Printf("order status updated: %v", order)
	}
}

func (a *Accrual) HandleStatusProcessing(ctx context.Context) {
	// Get all orders with status PROCESSING
	orders, err := a.order.GetByStatus(ctx, model.OrderStatusProcessing)
	if err != nil {
		log.Printf("error getting orders with status PROCESSING: %v", err)
		return
	}
	log.Printf("orders with status PROCESSING: %v", orders)
	// Get info from accrual service for all orders
	for _, order := range orders {
		accrualOrder, err := a.GetOrder(order.Number)
		if err != nil || accrualOrder.Status == OrderStatusProcessing || accrualOrder.Status == OrderStatusRegistered {
			log.Printf("error getting order from accrual service: %v", err)
			continue
		}
		log.Printf("order from accrual service: %v", accrualOrder)
		// Update order status and accrual in DB
		order.Status = MapOrderStatus(accrualOrder.Status)
		if accrualOrder.Status == OrderStatusProcessed && accrualOrder.Accrual != nil {
			order.Accrual = accrualOrder.Accrual
		}
		if err := a.order.Update(ctx, order); err != nil {
			log.Printf("error updating order: %v", err)
			continue
		}
		log.Printf("order updated: %v", order)
		// Update balance in DB
		if order.Status == model.OrderStatusProcessed {
			balance := model.Balance{
				UserLogin:   order.UserLogin,
				OrderNumber: order.Number,
				Debit:       order.Accrual,
				Credit:      nil,
				ProcessedAt: time.Now(),
			}
			if err := a.balance.Create(ctx, &balance); err != nil {
				log.Printf("error creating balance: %v", err)
				continue
			}
			log.Printf("balance created: %v", balance)
		}
	}
}

func (a *Accrual) GetOrder(orderNumber string) (*AccrualOrder, error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.address, orderNumber)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		defer response.Body.Close()
		payload, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var accrual AccrualOrder
		if err := json.Unmarshal(payload, &accrual); err != nil {
			return nil, err
		}
		return &accrual, nil
	case http.StatusNoContent:
	case http.StatusTooManyRequests:
	case http.StatusInternalServerError:
	}

	return nil, ErrInvalidResponseStatus
}

var ErrInvalidResponseStatus = errors.New("invalid response status")

type OrderStatus string

const (
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type AccrualOrder struct {
	Order   string
	Status  OrderStatus
	Accrual *decimal.Decimal
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
