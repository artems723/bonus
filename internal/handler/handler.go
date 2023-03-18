package handler

import (
	"bonus/internal/service"
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
		r.Get("/api/user/balance", h.TempHandler)
		r.Post("/api/user/balance/withdraw", h.TempHandler)
		r.Get("/api/user/withdrawals", h.TempHandler)
	})
	return r
}

func (h *Handler) TempHandler(w http.ResponseWriter, r *http.Request) {

}
