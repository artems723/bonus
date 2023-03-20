package handler

import (
	"bonus/internal/model"
	"bonus/internal/service"
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"net/http"
	"sort"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(LoginKey).(string)
	if !ok {
		http.Error(w, "no login in context", http.StatusInternalServerError)
		return
	}
	currentBalance, err := h.balanceService.GetByLogin(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Encode to JSON and write to response
	decimal.MarshalJSONWithoutQuotes = true
	err = json.NewEncoder(w).Encode(currentBalance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(LoginKey).(string)
	if !ok {
		http.Error(w, "no login in context", http.StatusInternalServerError)
		return
	}

	var withdrawal *model.Withdrawal
	// Read JSON and store to user struct
	err := json.NewDecoder(r.Body).Decode(&withdrawal)
	// Check errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//ok, err = h.orderService.CheckOrder(r.Context(), login, withdrawal.Order)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if !ok {
	//	http.Error(w, "wrong order number", http.StatusUnprocessableEntity)
	//	return
	//}

	err = h.balanceService.Withdraw(r.Context(), login, withdrawal)
	if err != nil && !errors.Is(err, service.ErrNotEnoughFunds) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, service.ErrNotEnoughFunds) {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(LoginKey).(string)
	if !ok {
		http.Error(w, "no login in context", http.StatusInternalServerError)
		return
	}
	withdrawals, err := h.balanceService.GetWithdrawals(r.Context(), login)
	if err != nil && !errors.Is(err, service.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].ProcessedAt.Before(withdrawals[j].ProcessedAt)
	})

	w.Header().Set("Content-Type", "application/json")
	// Encode to JSON and write to response
	err = json.NewEncoder(w).Encode(withdrawals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
