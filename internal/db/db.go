package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// Pool is a type alias for pgxpool.Pool to abstract the database driver.
type Pool = pgxpool.Pool

// NewPool creates a PostgreSQL connection pool with sensible defaults.
// Logs connection success with pool configuration details.
func NewPool(ctx context.Context, databaseURL string, log zerolog.Logger) (*Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("host", config.ConnConfig.Host).
		Uint16("port", config.ConnConfig.Port).
		Str("database", config.ConnConfig.Database).
		Int32("max_conns", config.MaxConns).
		Int32("min_conns", config.MinConns).
		Msg("Database connection established")

	return pool, nil
}

// Close gracefully closes the connection pool.
func Close(pool *Pool) {
	if pool != nil {
		pool.Close()
	}
}
