package grpc_api

import (
	"context"
	"errors"
	"testing"

	"github.com/Sugyk/auth_service/internal/api/grpc/pb"
	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/pkg/logger"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	s := NewServer(mockService, logger.NewNoop())

	tests := []struct {
		name         string
		req          *pb.RegisterRequest
		setupMock    func()
		expectedMsg  string
		expectedCode codes.Code
		expectedErr  string
	}{
		{
			name: "success - user successfully registered",
			req:  &pb.RegisterRequest{Login: "newuser", Password: "StrongPass12345678!"},
			setupMock: func() {
				mockService.EXPECT().
					Register(gomock.Any(), "newuser", "StrongPass12345678!").
					Return(nil)
			},
			expectedMsg: "user with login 'newuser' created",
		},
		{
			name:         "validation error - empty login",
			req:          &pb.RegisterRequest{Login: "", Password: "StrongPass12345678!"},
			setupMock:    func() {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "validation error - short password",
			req:          &pb.RegisterRequest{Login: "validuser", Password: "123"},
			setupMock:    func() {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "conflict - login already exists",
			req:  &pb.RegisterRequest{Login: "existinguser", Password: "StrongPass12345678!"},
			setupMock: func() {
				mockService.EXPECT().
					Register(gomock.Any(), "existinguser", "StrongPass12345678!").
					Return(models.NewDuplicateLoginErr("existinguser"))
			},
			expectedCode: codes.AlreadyExists,
			expectedErr:  "Login existinguser is already exists",
		},
		{
			name: "internal error from service - cause not leaked",
			req:  &pb.RegisterRequest{Login: "testuser", Password: "StrongPass12345678!"},
			setupMock: func() {
				mockService.EXPECT().
					Register(gomock.Any(), "testuser", "StrongPass12345678!").
					Return(errors.New("database timeout"))
			},
			expectedCode: codes.Internal,
			expectedErr:  "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := s.Register(context.Background(), tt.req)

			if tt.expectedCode == codes.OK {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp.GetMessage() != tt.expectedMsg {
					t.Errorf("expected message %q, got %q", tt.expectedMsg, resp.GetMessage())
				}
				return
			}

			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("expected a gRPC status error, got %v", err)
			}
			if st.Code() != tt.expectedCode {
				t.Errorf("expected code %s, got %s", tt.expectedCode, st.Code())
			}
			if tt.expectedErr != "" && st.Message() != tt.expectedErr {
				t.Errorf("expected message %q, got %q", tt.expectedErr, st.Message())
			}
		})
	}
}

func TestServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	s := NewServer(mockService, logger.NewNoop())

	tests := []struct {
		name         string
		req          *pb.LoginRequest
		setupMock    func()
		expectedTok  string
		expectedCode codes.Code
		expectedErr  string
	}{
		{
			name: "success - token issued",
			req:  &pb.LoginRequest{Login: "existing_user", Password: "StrongPass12345678!"},
			setupMock: func() {
				mockService.EXPECT().
					Login(gomock.Any(), "existing_user", "StrongPass12345678!").
					Return("test_token", nil)
			},
			expectedTok: "test_token",
		},
		{
			name:         "validation error - empty login",
			req:          &pb.LoginRequest{Login: "", Password: "StrongPass12345678!"},
			setupMock:    func() {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "validation error - short password",
			req:          &pb.LoginRequest{Login: "validuser", Password: "123"},
			setupMock:    func() {},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "internal error from service - cause not leaked",
			req:  &pb.LoginRequest{Login: "existinguser", Password: "StrongPass12345678!"},
			setupMock: func() {
				mockService.EXPECT().
					Login(gomock.Any(), "existinguser", "StrongPass12345678!").
					Return("", errors.New("database timeout"))
			},
			expectedCode: codes.Internal,
			expectedErr:  "internal error",
		},
		{
			name: "login not found",
			req:  &pb.LoginRequest{Login: "unexistinguser", Password: "StrongPass12345678!"},
			setupMock: func() {
				mockService.EXPECT().
					Login(gomock.Any(), "unexistinguser", "StrongPass12345678!").
					Return("", models.NewLoginNotFound("unexistinguser"))
			},
			expectedCode: codes.Unauthenticated,
			expectedErr:  "There is incorrect login or password",
		},
		{
			name: "wrong password",
			req:  &pb.LoginRequest{Login: "existinguser", Password: "WrongStrongPass12345678!"},
			setupMock: func() {
				mockService.EXPECT().
					Login(gomock.Any(), "existinguser", "WrongStrongPass12345678!").
					Return("", models.NewWrongPassword("existinguser"))
			},
			expectedCode: codes.Unauthenticated,
			expectedErr:  "There is incorrect login or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := s.Login(context.Background(), tt.req)

			if tt.expectedCode == codes.OK {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp.GetAccessToken() != tt.expectedTok {
					t.Errorf("expected token %q, got %q", tt.expectedTok, resp.GetAccessToken())
				}
				return
			}

			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("expected a gRPC status error, got %v", err)
			}
			if st.Code() != tt.expectedCode {
				t.Errorf("expected code %s, got %s", tt.expectedCode, st.Code())
			}
			if tt.expectedErr != "" && st.Message() != tt.expectedErr {
				t.Errorf("expected message %q, got %q", tt.expectedErr, st.Message())
			}
		})
	}
}

func TestToGRPCError_UnknownCode(t *testing.T) {
	var unknownCode models.ErrorCode = "test_code"
	appErr := models.New(unknownCode, "details", "cause")

	err := toGRPCError(appErr)

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected a gRPC status error, got %v", err)
	}
	if st.Code() != codes.Internal {
		t.Errorf("expected fallback code %s, got %s", codes.Internal, st.Code())
	}
}
