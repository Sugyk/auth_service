package service

import (
	"context"
	"fmt"
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
}

func NewService(repo Repository, txManager TxManager) *Service {
	return &Service{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *Service) Register(ctx context.Context, login string, password string) error {
	txErr := s.txManager.WithTx(
		ctx,
		func(ctx context.Context) error {
			if err := s.repo.CreateUser(ctx, login, password); err != nil {
				return fmt.Errorf("creating user: %w", err)
			}
			return nil
		},
	)
	if txErr != nil {
		return fmt.Errorf("error executing tx: %w", txErr)
	}

	return nil
}
