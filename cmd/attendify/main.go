// Package main is the entry point for the Attendify application.
// Attendify is a real-time attendance backend system that enables:
//   - Teachers to create and manage classes
//   - Students to mark attendance in real-time via WebSocket
//   - Persistent storage of attendance records in PostgreSQL
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/tahiriqbal095/attendify/internal/config"
	"github.com/tahiriqbal095/attendify/internal/db"
	"github.com/tahiriqbal095/attendify/internal/logger"
	"github.com/tahiriqbal095/attendify/internal/server"
)

const (
	// dbConnectTimeout is the maximum time allowed to establish a database connection.
	dbConnectTimeout = 10 * time.Second

	// shutdownTimeout is the maximum time allowed for graceful shutdown.
	// During this period, the server stops accepting new connections
	// and waits for existing requests to complete.
	shutdownTimeout = 5 * time.Second
)

func main() {
	// Load configuration from environment variables.
	// Exits immediately if required config is missing.
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize structured logger based on environment (development/production).
	log := logger.NewLogger(cfg.Environment)

	// Establish database connection pool.
	// Uses a timeout to prevent hanging on unreachable database.
	pool, err := connectDatabase(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close(pool)

	// Create and start HTTP server.
	srv := server.NewServer(cfg.AppPort, log, pool)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Block until shutdown signal is received.
	sig := waitForShutdownSignal()
	log.Info().Str("signal", sig.String()).Msg("Shutdown signal received")

	// Gracefully shutdown the server.
	if err := gracefulShutdown(srv); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Server exited gracefully")
}

// connectDatabase establishes a connection pool to PostgreSQL.
// Returns an error if the connection cannot be established within the timeout.
func connectDatabase(databaseURL string, log zerolog.Logger) (*db.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
	defer cancel()

	return db.NewPool(ctx, databaseURL, log)
}

// waitForShutdownSignal blocks until SIGINT or SIGTERM is received.
// SIGINT is triggered by Ctrl+C, SIGTERM by `kill` or container orchestrators.
func waitForShutdownSignal() os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return <-quit
}

// gracefulShutdown stops the server gracefully within the timeout period.
// It stops accepting new connections and waits for active requests to complete.
func gracefulShutdown(srv *server.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctx)
}
