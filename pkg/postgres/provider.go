package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
}

type Provider struct {
	pool   *pgxpool.Pool
	logger Logger

	ConnStr         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime int
	MaxConnIdleTime int
}

func (p *Provider) Open(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(p.ConnStr)
	if err != nil {
		return fmt.Errorf("parse connection string: %w", err)
	}

	poolConfig.MaxConns = p.MaxConns
	poolConfig.MinConns = p.MinConns
	poolConfig.MaxConnLifetime = time.Duration(p.MaxConnLifetime) * time.Minute
	poolConfig.MaxConnIdleTime = time.Duration(p.MaxConnIdleTime) * time.Minute

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
		ctx,
		"Connected to PostgreSQL",
		"host", poolConfig.ConnConfig.Host,
		"port", poolConfig.ConnConfig.Port,
	)

	return nil
}

func (p *Provider) DB() *pgxpool.Pool {
	return p.pool
}

func (p *Provider) Close() {
	p.pool.Close()
}

func NewProvider(logger Logger,
	ConnStr string,
	MaxConns int32,
	MinConns int32,
	MaxConnLifetime int,
	MaxConnIdleTime int,
) *Provider {
	return &Provider{
		logger:          logger,
		ConnStr:         ConnStr,
		MaxConns:        MaxConns,
		MinConns:        MinConns,
		MaxConnLifetime: MaxConnLifetime,
		MaxConnIdleTime: MaxConnIdleTime,
	}
}
