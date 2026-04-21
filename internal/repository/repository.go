package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateUser(ctx context.Context, login string, password string) error {
	exec := postgres.GetExecutor(ctx, r.pool)

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

func (r *Repository) GetPasswordByLogin(ctx context.Context, login string) (string, error) {
	exec := postgres.GetExecutor(ctx, r.pool)

	query := `
		SELECT * FROM Users
		WHERE login = $1
	`

	row := exec.QueryRow(ctx, query, login)

	var passHash string

	err := row.Scan(&passHash)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", models.ErrLoginNotFound
	}

	return passHash, nil
}
