package handler

import (
	"bonus/internal/model"
	"bonus/internal/repository"
	"bonus/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
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
		r.Post("/api/user/orders", h.TempHandler)
		r.Get("/api/user/orders", h.TempHandler)
		r.Get("/api/user/balance", h.TempHandler)
		r.Post("/api/user/balance/withdraw", h.TempHandler)
		r.Get("/api/user/withdrawals", h.TempHandler)
	})
	return r
}

func (h *Handler) TempHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("hi")))
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	// Read JSON and store to user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	// Check errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userService.Create(r.Context(), user)
	if err != nil && !errors.Is(err, repository.ErrUsernameIsTaken) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, repository.ErrUsernameIsTaken) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	log.Printf("User was registered: %v\n", user)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("please login")))
}

// Authenticator middleware
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		user := model.User{
			Login:        username,
			PasswordHash: password,
		}
		log.Printf("User: %v\n", user)
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		// authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
