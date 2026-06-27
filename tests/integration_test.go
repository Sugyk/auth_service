package integration_tests

import (
	"net/http"

	"github.com/Sugyk/auth_service/internal/models"
)

const ServicePrefix = "/api/v1/auth"

func (s *IntegrationSuite) TestRegister() {
	tests := []struct {
		name           string
		req            models.RegisterRequest
		expectedStatus int
	}{
		// ✅ Happy path
		{
			name:           "valid: minimal password 16 chars",
			req:            models.RegisterRequest{Login: "testuser", Password: "1234567890123456"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid: long password",
			req:            models.RegisterRequest{Login: "testuser2", Password: "averylongpasswordthatisvalid123"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid: password exactly 16 unicode runes",
			req:            models.RegisterRequest{Login: "testuser3", Password: "пароль123456789!"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid: login can be anything non-empty",
			req:            models.RegisterRequest{Login: "u", Password: "1234567890123456"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid: login with spaces",
			req:            models.RegisterRequest{Login: "test user", Password: "1234567890123456"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid: login with special chars",
			req:            models.RegisterRequest{Login: "!@#$%^&*()", Password: "1234567890123456"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid: login with emoji",
			req:            models.RegisterRequest{Login: "😀", Password: "1234567890123456"},
			expectedStatus: http.StatusCreated,
		},
		// ❌ Empty login
		{
			name:           "invalid: empty login",
			req:            models.RegisterRequest{Login: "", Password: "1234567890123456"},
			expectedStatus: http.StatusBadRequest,
		},
		// ❌ Password too short (< 16 runes)
		{
			name:           "invalid: password 15 chars",
			req:            models.RegisterRequest{Login: "userA", Password: "123456789012345"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid: password 1 char",
			req:            models.RegisterRequest{Login: "userB", Password: "x"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid: empty password",
			req:            models.RegisterRequest{Login: "userC", Password: ""},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid: password 15 unicode runes",
			req:            models.RegisterRequest{Login: "userD", Password: "пароль1234567!"},
			expectedStatus: http.StatusBadRequest,
		},

		// ❌ Both invalid
		{
			name:           "invalid: empty login and short password",
			req:            models.RegisterRequest{Login: "", Password: "short"},
			expectedStatus: http.StatusBadRequest,
		},

		// ❌ Duplicate login
		{
			name:           "duplicate: first registration",
			req:            models.RegisterRequest{Login: "duplicate_user", Password: "1234567890123456"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "duplicate: second registration same login",
			req:            models.RegisterRequest{Login: "duplicate_user", Password: "differentpassword1234"},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			resp := s.PerformRequest(
				http.MethodPost,
				ServicePrefix+"/reg",
				tc.req,
				s.handler.Register,
			)
			s.Equal(tc.expectedStatus, resp.Code, "test case: %s\nresponse body: %s", tc.name, resp.Body.String())
		})
	}
}
