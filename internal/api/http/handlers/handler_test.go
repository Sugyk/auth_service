package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/pkg/logger"
	"go.uber.org/mock/gomock"
)

func TestHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noopLogger := logger.NewNoop()

	mockService := NewMockService(ctrl)
	h := NewHandler(mockService, noopLogger)

	tests := []struct {
		name           string
		requestBody    any
		setupMock      func()
		expectedStatus int
		expectedErr    *models.AppError
	}{
		{
			name: "success - user successfully registered",
			requestBody: models.RegisterRequest{
				Login:    "newuser",
				Password: "StrongPass12345678!",
			},
			setupMock: func() {
				mockService.EXPECT().
					Register(gomock.Any(), "newuser", "StrongPass12345678!").
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "bad request - invalid JSON",
			requestBody:    `{"login": "user", "password": "pass"`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    models.NewValidationErr("error decoding body"),
		},
		{
			name: "validation error - empty login",
			requestBody: models.RegisterRequest{
				Login:    "",
				Password: "StrongPass12345678!",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    models.NewValidationErr(""),
		},
		{
			name: "validation error - short password",
			requestBody: models.RegisterRequest{
				Login:    "validuser",
				Password: "123",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    models.NewValidationErr(""),
		},
		{
			name: "conflict - login already exists",
			requestBody: models.RegisterRequest{
				Login:    "existinguser",
				Password: "StrongPass12345678!",
			},
			setupMock: func() {
				mockService.EXPECT().
					Register(gomock.Any(), "existinguser", "StrongPass12345678!").
					Return(models.NewDuplicateLoginErr("existinguser"))
			},
			expectedStatus: http.StatusConflict,
			expectedErr:    models.NewDuplicateLoginErr("existinguser"),
		},
		{
			name: "internal server error from service",
			requestBody: models.RegisterRequest{
				Login:    "testuser",
				Password: "StrongPass12345678!",
			},
			setupMock: func() {
				mockService.EXPECT().
					Register(gomock.Any(), "testuser", "StrongPass12345678!").
					Return(errors.New("database timeout"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    models.NewInternalErr(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var body []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			h.Register(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedErr != nil {
				var gotErr models.AppError
				if err := json.NewDecoder(rr.Body).Decode(&gotErr); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}

				if gotErr.ErrCode != tt.expectedErr.ErrCode {
					t.Errorf("expected error code %q, got %q", tt.expectedErr.ErrCode, gotErr.ErrCode)
				}

				if tt.expectedErr.ErrCode != models.CodeValidationError {
					if gotErr.Details != tt.expectedErr.Details {
						t.Errorf("expected details %q, got %q", tt.expectedErr.Details, gotErr.Details)
					}
				}
			} else {
				var resp models.RegisterResponse201
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode success response: %v", err)
				}

				expectedMsg := "user with login 'newuser' created"
				if resp.Message != expectedMsg {
					t.Errorf("expected message %q, got %q", expectedMsg, resp.Message)
				}
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	noopLogger := logger.NewNoop()

	mockService := NewMockService(ctrl)
	h := NewHandler(mockService, noopLogger)

	tests := []struct {
		name           string
		requestBody    any
		setupMock      func()
		expectedStatus int
		expectedErr    *models.AppError
		expectedBody   any
	}{
		{
			name: "success - user successfully registered",
			requestBody: models.LoginRequest{
				Login:    "existing_user",
				Password: "StrongPass12345678!",
			},
			setupMock: func() {
				mockService.EXPECT().
					Login(gomock.Any(), "existing_user", "StrongPass12345678!").
					Return("test_token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   models.LoginResponse200{AccessToken: "test_token"},
		},
		{
			name:           "bad request - invalid JSON",
			requestBody:    `{"login": "user", "password": "pass"`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    models.NewValidationErr("error decoding body"),
		},
		{
			name: "validation error - empty login",
			requestBody: models.LoginRequest{
				Login:    "",
				Password: "StrongPass12345678!",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    models.NewValidationErr(""),
		},
		{
			name: "validation error - short password",
			requestBody: models.LoginRequest{
				Login:    "validuser",
				Password: "123",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    models.NewValidationErr(""),
		},
		{
			name: "internal server error from service",
			requestBody: models.LoginRequest{
				Login:    "existinguser",
				Password: "StrongPass12345678!",
			},
			setupMock: func() {
				mockService.EXPECT().
					Login(gomock.Any(), "existinguser", "StrongPass12345678!").
					Return("", errors.New("database timeout"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    models.NewInternalErr(),
		},
		{
			name: "Login not found",
			requestBody: models.LoginRequest{
				Login:    "unexistinguser",
				Password: "StrongPass12345678!",
			},
			setupMock: func() {
				mockService.EXPECT().
					Login(gomock.Any(), "unexistinguser", "StrongPass12345678!").
					Return("", models.NewLoginNotFound())
			},
			expectedStatus: http.StatusUnauthorized,
			expectedErr:    models.NewLoginNotFound(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var body []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			h.Login(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if tt.expectedBody != nil {

			}
			if tt.expectedErr != nil {
				var gotErr models.AppError
				if err := json.NewDecoder(rr.Body).Decode(&gotErr); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}

				if gotErr.ErrCode != tt.expectedErr.ErrCode {
					t.Errorf("expected error code %q, got %q", tt.expectedErr.ErrCode, gotErr.ErrCode)
				}

				if tt.expectedErr.ErrCode != models.CodeValidationError {
					if gotErr.Details != tt.expectedErr.Details {
						t.Errorf("expected details %q, got %q", tt.expectedErr.Details, gotErr.Details)
					}
				}
			} else {
				var resp models.LoginResponse200
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode success response: %v", err)
				}

				if resp != tt.expectedBody {
					t.Errorf("expected message %q, got %q", resp, tt.expectedBody)
				}
			}
		})
	}

}
