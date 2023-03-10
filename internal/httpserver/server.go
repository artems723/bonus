package httpserver

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New() Server {
	return Server{}
}

func (s *Server) Run(ctx context.Context, serverAddr string, r *chi.Mux) error {
	s.httpServer = &http.Server{
		Addr:           serverAddr,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        r,
	}

	// handle graceful termination
	go func() {
		<-ctx.Done()
		if s.httpServer != nil {
			err := s.httpServer.Shutdown(ctx)
			if err != nil {
				log.Printf("failed to shut down http server, %v", err)
			}
		}
	}()

	return s.httpServer.ListenAndServe()
}
