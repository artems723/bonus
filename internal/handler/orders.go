package handler

import (
	"bonus/internal/model"
	"bonus/internal/service"
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"time"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(LoginKey).(string)
	if !ok {
		http.Error(w, "no login in context", http.StatusInternalServerError)
		return
	}

	orderNumber, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order := model.Order{
		UserLogin:  login,
		Number:     string(orderNumber),
		Status:     model.OrderStatusNew,
		UploadedAt: time.Now(),
	}

	if !order.Valid() {
		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	err = h.orderService.Create(r.Context(), &order)
	if err != nil && !errors.Is(err, service.ErrOrderAlreadyExists) && !errors.Is(err, service.ErrOrderAlreadyExistsForAnotherUser) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, service.ErrOrderAlreadyExists) {
		w.Write([]byte("order already exists"))
		w.WriteHeader(http.StatusOK)
		return
	}
	if errors.Is(err, service.ErrOrderAlreadyExistsForAnotherUser) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(LoginKey).(string)
	if !ok {
		http.Error(w, "no login in context", http.StatusInternalServerError)
		return
	}
	orders, err := h.orderService.GetByLogin(r.Context(), login)
	if err != nil && !errors.Is(err, service.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Encode to JSON and write to response
	decimal.MarshalJSONWithoutQuotes = true
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
