package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/pyr33x/ev/internal/server/http"
	"github.com/pyr33x/ev/pkg/config"

	h "net/http"
)

func main() {
	config := config.New()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := http.New(&config.Server)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != h.ErrServerClosed {
			log.Println("failed to listen and serve", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("http server shutdown error", err)
	}

	log.Println("server stopped gracefully.")
}
