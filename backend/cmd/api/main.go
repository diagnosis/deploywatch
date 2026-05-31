// backend/cmd/api/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diagnosis/deploywatchv2/internal/application"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := application.NewApplication()
	if err != nil {
		logger.Error(ctx, "failed start application", "err", err)
		os.Exit(1)
	}

	r := application.SetupRoutes(app)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", "8080"),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second, // only timeout headers, not body (SSE-safe)
		WriteTimeout:      0,                // no write timeout — required for SSE
		IdleTimeout:       120 * time.Second,
	}

	app.StartEventCleanup()
	go app.Hub.Run(ctx)
	go app.WsHub.Run(ctx)

	go func() {
		logger.Info(ctx, "server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, "server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info(ctx, "shutdown signal received", "signal", sig)

	// Give in-flight requests 30 seconds to complete
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "graceful shutdown failed", "err", err)
		os.Exit(1)
	}

	logger.Info(ctx, "server stopped gracefully")
}
