package repository

import (
	"context"
	"database/sql"
	"fmt"

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
        RETURNING id
    `

	var id int
	err := exec.QueryRow(ctx, query, login, password).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("login '%s' is already exists", login)
		}
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}
