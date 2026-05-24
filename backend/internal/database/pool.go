package database

import (
	"context"

	"github.com/diagnosis/deploywatchv2/internal/config"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenPool(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	pgxConfig, err := pgxpool.ParseConfig(cfg.Database.Dsn)
	if err != nil {
		return nil, err
	}
	logger.Info(ctx, "connecting to db...")

	pgxConfig.MinConns = cfg.Database.MinConns
	pgxConfig.MaxConns = cfg.Database.MaxConns
	pgxConfig.MaxConnLifetime = cfg.Database.MaxConnLifetime
	pgxConfig.MaxConnIdleTime = cfg.Database.MaxConnIdleTime
	pgxConfig.HealthCheckPeriod = cfg.Database.HealthCheckPeriod
	pgxConfig.ConnConfig.ConnectTimeout = cfg.Database.ConnectTimeout

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	logger.Info(ctx, "connected to db")
	return pool, nil
}
