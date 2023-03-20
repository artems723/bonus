package handler

import (
	"bonus/internal/model"
	"bonus/internal/service"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Handler struct {
	userService    service.UserService
	orderService   service.OrderService
	balanceService service.BalanceService
}

func New(u service.UserService, o service.OrderService, b service.BalanceService) *Handler {
	return &Handler{
		userService:    u,
		orderService:   o,
		balanceService: b,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	// Create new chi router
	r := chi.NewRouter()

	// Using built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.Compress(5))
	r.Use(middleware.Recoverer)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", h.RegisterHandler)
		r.Post("/api/user/login", h.LoginHandler)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(Authenticator)
		r.Post("/api/user/orders", h.CreateOrder)
		r.Get("/api/user/orders", h.GetOrders)
		r.Get("/api/user/balance", h.GetBalance)
		r.Post("/api/user/balance/withdraw", h.Withdraw)
		r.Get("/api/user/withdrawals", h.TempHandler)
	})
	return r
}

func (h *Handler) TempHandler(w http.ResponseWriter, r *http.Request) {

}

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

	ok, err = h.orderService.CheckOrder(r.Context(), login, withdrawal.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "wrong order number", http.StatusUnprocessableEntity)
		return
	}

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
