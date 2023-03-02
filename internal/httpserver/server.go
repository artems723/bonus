package httpserver

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type server struct {
	httpServer *http.Server
}

func New() server {
	return server{}
}

func (s *server) Run(serverAddr string, r *chi.Mux) error {
	s.httpServer = &http.Server{
		Addr:           serverAddr,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        r,
	}
	return s.httpServer.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
