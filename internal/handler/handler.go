package handler

import (
	"bonus/internal/service"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"sort"
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
		r.Get("/api/user/withdrawals", h.Withdrawals)
	})
	return r
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
