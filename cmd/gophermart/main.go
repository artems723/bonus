package main

import (
	"bonus/internal/config"
	"bonus/internal/handler"
	"bonus/internal/httpserver"
	"bonus/internal/repository"
	"bonus/internal/service"
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create and read config
	cfg := config.Config{}
	// Parse config from flag
	flag.StringVar(&cfg.RunAddress, "a", ":8080", "server address and port")
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	flag.StringVar(&cfg.DatabaseURI, "d", "postgres://postgres:pass@postgres/postgres?sslmode=disable", "database Uniform Resource Identifier")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "accrual system address")
	flag.Parse()
	// Parse config from env
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("error parsing config file: %v", err)
	}
	log.Printf("Using config: %#v", cfg)

	// Connect to DB
	db, err := sqlx.Connect("pgx", cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}

	// Wait for successful DB ping
	for {
		err := db.Ping()
		if err != nil {
			log.Printf("db ping: %v\n", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	// Apply migrations

	log.Println("connected to DB")

	// Create repositories
	orderRepo := repository.NewOrderRepository(db)
	balanceRepo := repository.NewBalanceRepository(db)
	// Create service
	serv := service.New(orderRepo, balanceRepo)

	// Create handler
	h := handler.New(serv)

	// Create server
	srv := httpserver.New()

	// Create channel for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start http server
	go func() {
		err = srv.Run(cfg.RunAddress, h.InitRoutes())
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("srv.Run, error occured while running http server: %v", err)
		}
	}()
	log.Printf("Server started")

	<-done
	// Shutdown http server
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("Server shutdown Failed:%+v", err)
	}
	err = serv.Shutdown()
	if err != nil {
		log.Fatalf("serv.Shutdown: %v", err)
	}
	log.Print("Server stopped properly")

}
