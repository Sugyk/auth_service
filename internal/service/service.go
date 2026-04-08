package service

import "context"

type Repository interface{}

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
