package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Sugyk/auth_service/internal/models"
	"go.uber.org/mock/gomock"
)

func withTxRunFn(ctx context.Context) func(context.Context, func(context.Context) error) error {
	return func(_ context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	}
}

func TestService_Register(t *testing.T) {
	tests := []struct {
		name      string
		login     string
		password  string
		setupMock func(repo *MockRepository, tx *MockTxManager, hasher *MockPasswordHasher, jwtMgr *MockJWTManager)
		wantErr   error
		wantErrAs bool
	}{
		{
			name:     "success",
			login:    "john",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, tx *MockTxManager, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				hasher.EXPECT().HashPassword("StrongPass12345678!").Return("hashed", nil)
				tx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(withTxRunFn(context.Background()))
				repo.EXPECT().CreateUser(gomock.Any(), "john", "hashed").Return(nil)
			},
		},
		{
			name:     "hasher error - tx never started",
			login:    "john",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, tx *MockTxManager, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				hasher.EXPECT().HashPassword("StrongPass12345678!").Return("", errors.New("bcrypt exploded"))
			},
			wantErr: errors.New("hashing password: bcrypt exploded"),
		},
		{
			name:     "duplicate login is translated to AppError",
			login:    "existing",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, tx *MockTxManager, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				hasher.EXPECT().HashPassword("StrongPass12345678!").Return("hashed", nil)
				tx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(withTxRunFn(context.Background()))
				repo.EXPECT().CreateUser(gomock.Any(), "existing", "hashed").Return(models.ErrDuplicate)
			},
			wantErrAs: true,
			wantErr:   models.NewDuplicateLoginErr("existing"),
		},
		{
			name:     "other repository error is wrapped, not swallowed",
			login:    "john",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, tx *MockTxManager, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				hasher.EXPECT().HashPassword("StrongPass12345678!").Return("hashed", nil)
				tx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(withTxRunFn(context.Background()))
				repo.EXPECT().CreateUser(gomock.Any(), "john", "hashed").Return(errors.New("connection reset"))
			},
			wantErr: errors.New("creating user: connection reset"),
		},
		{
			name:     "tx manager itself fails - repo never reached",
			login:    "john",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, tx *MockTxManager, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				hasher.EXPECT().HashPassword("StrongPass12345678!").Return("hashed", nil)
				tx.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(errors.New("begin transaction: connection refused"))
			},
			wantErr: errors.New("begin transaction: connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := NewMockRepository(ctrl)
			tx := NewMockTxManager(ctrl)
			hasher := NewMockPasswordHasher(ctrl)
			jwtMgr := NewMockJWTManager(ctrl)
			tt.setupMock(repo, tx, hasher, jwtMgr)

			svc := NewService(repo, tx, hasher, jwtMgr)

			err := svc.Register(context.Background(), tt.login, tt.password)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}

			if tt.wantErrAs {
				var appErr *models.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("expected an *models.AppError, got %T: %v", err, err)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %q, got %q", tt.wantErr, err)
				}
				return
			}

			if err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %q, got %q", tt.wantErr, err)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	tests := []struct {
		name      string
		login     string
		password  string
		setupMock func(repo *MockRepository, hasher *MockPasswordHasher, jwtMgr *MockJWTManager)
		wantToken string
		wantErr   error
		wantErrAs bool
	}{
		{
			name:     "success",
			login:    "john",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				repo.EXPECT().GetPasswordByLogin(gomock.Any(), "john").Return("hashed", nil)
				hasher.EXPECT().CompareHashAndPassword("StrongPass12345678!", "hashed").Return(true)
				jwtMgr.EXPECT().CreateJWT("john").Return("token123", nil)
			},
			wantToken: "token123",
		},
		{
			name:     "login not found - generic AppError, not leaked",
			login:    "ghost",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				repo.EXPECT().GetPasswordByLogin(gomock.Any(), "ghost").Return("", models.ErrLoginNotFound)
			},
			wantErrAs: true,
			wantErr:   models.NewLoginNotFound("ghost"),
		},
		{
			name:     "other repository error is wrapped",
			login:    "john",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				repo.EXPECT().GetPasswordByLogin(gomock.Any(), "john").Return("", errors.New("connection reset"))
			},
			wantErr: errors.New("get password by login: connection reset"),
		},
		{
			name:     "wrong password - same generic AppError as not-found",
			login:    "john",
			password: "WrongPass12345678!",
			setupMock: func(repo *MockRepository, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				repo.EXPECT().GetPasswordByLogin(gomock.Any(), "john").Return("hashed", nil)
				hasher.EXPECT().CompareHashAndPassword("WrongPass12345678!", "hashed").Return(false)
			},
			wantErrAs: true,
			wantErr:   models.NewWrongPassword("john"),
		},
		{
			name:     "jwt signing failure becomes internal AppError",
			login:    "john",
			password: "StrongPass12345678!",
			setupMock: func(repo *MockRepository, hasher *MockPasswordHasher, jwtMgr *MockJWTManager) {
				repo.EXPECT().GetPasswordByLogin(gomock.Any(), "john").Return("hashed", nil)
				hasher.EXPECT().CompareHashAndPassword("StrongPass12345678!", "hashed").Return(true)
				jwtMgr.EXPECT().CreateJWT("john").Return("", errors.New("signing key missing"))
			},
			wantErrAs: true,
			wantErr:   models.NewInternalErr("signing key missing"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := NewMockRepository(ctrl)
			tx := NewMockTxManager(ctrl)
			hasher := NewMockPasswordHasher(ctrl)
			jwtMgr := NewMockJWTManager(ctrl)
			tt.setupMock(repo, hasher, jwtMgr)

			svc := NewService(repo, tx, hasher, jwtMgr)

			token, err := svc.Login(context.Background(), tt.login, tt.password)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if token != tt.wantToken {
					t.Errorf("expected token %q, got %q", tt.wantToken, token)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}
			if token != "" {
				t.Errorf("expected empty token on error, got %q", token)
			}

			if tt.wantErrAs {
				var appErr *models.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("expected an *models.AppError, got %T: %v", err, err)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %q, got %q", tt.wantErr, err)
				}
				return
			}

			if err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %q, got %q", tt.wantErr, err)
			}
		})
	}
}
