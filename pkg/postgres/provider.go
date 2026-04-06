package pgprovider

import (
	"context"
	"fmt"
	"time"

	"github.com/Sugyk/auth_service/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Logger interface {
	Info(msg string, args ...any)
}

type Provider struct {
	pool   *pgxpool.Pool
	logger Logger
	cfg    config.PgConfig
}

func (p *Provider) Open(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(p.cfg.ConnStr)
	if err != nil {
		return fmt.Errorf("parse connection string: %w", err)
	}

	poolConfig.MaxConns = p.cfg.MaxConns
	poolConfig.MinConns = p.cfg.MinConns
	poolConfig.MaxConnLifetime = time.Duration(p.cfg.MaxConnLifetime) * time.Minute
	poolConfig.MaxConnIdleTime = time.Duration(p.cfg.MaxConnIdleTime) * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return fmt.Errorf("ping database: %w", err)
	}

	p.pool = pool

	p.logger.Info(
		"Connected to PostgreSQL: %s/%s",
		poolConfig.ConnConfig.Host,
		poolConfig.ConnConfig.Port,
	)

	return nil
}

func NewProvider(logger Logger, config config.PgConfig) *Provider {
	return &Provider{
		logger: logger,
		cfg:    config,
	}
}
