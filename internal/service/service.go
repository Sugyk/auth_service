package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/internal/pkg/hasher"
)

type Repository interface {
	CreateUser(ctx context.Context, login string, password string) error
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	txManager TxManager
	repo      Repository
	hasher    hasher.PasswordHasher
}

func NewService(repo Repository, txManager TxManager, hasher hasher.PasswordHasher) *Service {
	return &Service{
		txManager: txManager,
		repo:      repo,
		hasher:    hasher,
	}
}

func (s *Service) Register(ctx context.Context, login string, password string) error {
	hashedPassword, err := s.hasher.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	if txErr := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		// Tx body start
		if err := s.repo.CreateUser(ctx, login, hashedPassword); err != nil {
			if errors.Is(err, models.ErrDuplicate) {
				return models.NewDuplicateLoginErr(login)
			}
			return fmt.Errorf("creating user: %w", err)
		}

		return nil
		// Tx body end
	}); txErr != nil {
		return txErr
	}

	return nil
}
