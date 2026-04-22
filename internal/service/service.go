package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/internal/pkg/hasher"
	"github.com/Sugyk/auth_service/internal/pkg/jwt_manager"
)

type Repository interface {
	CreateUser(ctx context.Context, login string, password string) error
	GetPasswordByLogin(ctx context.Context, login string) (string, error)
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	txManager  TxManager
	repo       Repository
	hasher     hasher.PasswordHasher
	jwtManager jwt_manager.JWTManager
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

// Login method boilerplate
func (s *Service) Login(ctx context.Context, login string, password string) (string, error) {

	passHash, err := s.repo.GetPasswordByLogin(ctx, login)

	if err != nil {
		if errors.Is(err, models.ErrLoginNotFound) {
			return "", models.NewLoginNotFound()
		}
		return "", fmt.Errorf("get password by login: %w", err)
	}

	if !s.hasher.CompareHashAndPassword(password, passHash) {
		return "", models.NewWrongPassword()
	}

	token, err := s.jwtManager.CreateJWT(login)

	if err != nil {
		return "", models.NewInternalErr()
	}

	return token, nil
}
