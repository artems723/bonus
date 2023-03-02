package handler

import (
	"bonus/internal/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", TempHandler)
		r.Post("/login", TempHandler)
		r.Post("/orders", TempHandler)
		r.Get("/orders", TempHandler)
		r.Get("/balance", TempHandler)
		r.Post("/balance/withdraw", TempHandler)
		r.Get("/withdrawals", TempHandler)
	})
	return r
}

func TempHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("hi")))
}
