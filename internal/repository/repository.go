package repository

import (
	"context"
	"fmt"

	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository() *Repository {
	return &Repository{}
}

func (d *Repository) CreateUser(ctx context.Context, login string, password string) error {
	exec := postgres.GetExecutor(ctx, d.pool)

	query := `
        INSERT INTO users (login, password_hash) 
        VALUES ($1, $2)
        ON CONFLICT (login) DO NOTHING
    `

	tag, err := exec.Exec(ctx, query, login, password)
	if err != nil {
		return fmt.Errorf("create user in postgres: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrDuplicate
	}
	return nil
}
