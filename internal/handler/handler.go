package handler

import (
	"bonus/internal/model"
	"bonus/internal/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type handler struct {
	s service.Service
}

func New(s service.Service) *handler {
	return &handler{
		s: s,
	}
}

func (h *handler) InitRoutes() *chi.Mux {
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
		r.Post("/api/user/register", RegisterHandler)
		r.Post("/api/user/login", LoginHandler)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(Authenticator)
		r.Post("/api/user/orders", TempHandler)
		r.Get("/api/user/orders", TempHandler)
		r.Get("/api/user/balance", TempHandler)
		r.Post("/api/user/balance/withdraw", TempHandler)
		r.Get("/api/user/withdrawals", TempHandler)
	})
	return r
}

func TempHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("hi")))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Unauthorized", http.StatusBadRequest)
	}
	user := model.User{
		Login:        username,
		PasswordHash: password,
	}

	w.Write([]byte(fmt.Sprintf("please register")))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
