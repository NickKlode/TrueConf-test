package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"refactoring/internal/api"
	"refactoring/internal/config"
	"refactoring/internal/storage"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	db := storage.New(cfg.Store)
	api := api.New(db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	srv := http.Server{
		Addr:         cfg.Address,
		Handler:      api.Router(),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server. %v", err)
		}
	}()
	log.Println("Server started...")

	<-quit
	log.Println("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed ti stop server. %v", err)
		return
	}
	log.Println("Server stopped")
}
