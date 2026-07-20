package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
}

type Provider struct {
	client *redis.Client
	logger Logger

	Addr     string
	Password string
	DB       int
}

func (p *Provider) Open(ctx context.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     p.Addr,
		Password: p.Password,
		DB:       p.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()

		return fmt.Errorf("ping redis: %w", err)
	}

	p.client = client

	p.logger.Info(ctx, "Connected to Redis", "addr", p.Addr, "db", p.DB)

	return nil
}

func (p *Provider) Client() *redis.Client {
	return p.client
}

func (p *Provider) Close() {
	p.client.Close()
}

func NewProvider(logger Logger, addr string, password string, db int) *Provider {
	return &Provider{
		logger:   logger,
		Addr:     addr,
		Password: password,
		DB:       db,
	}
}
