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

// TestLogin only covers the wrong-credentials path. A happy-path assertion
// (register then log in with the correct password) can't be exercised here:
// Register runs inside the suite's fixed, never-committed test transaction
// (see NewTestTxManager in pkg/postgres/tx.go), but Login's read
// (GetPasswordByLogin) is not wrapped in a transaction at all, so it always
// queries the real pool via GetExecutor and never sees the uncommitted row.
// This is the transaction test double limitation already called out in
// CLAUDE.md ("multi-transaction rollback scenarios can't be verified through
// it as-is") — fixing it would mean changing Login's transaction semantics,
// which is out of scope here.
func (s *IntegrationSuite) TestLogin() {
	resp := s.PerformRequest(
		http.MethodPost,
		ServicePrefix+"/login",
		models.LoginRequest{Login: "no_such_user", Password: "wrongpassword12345"},
		s.handler.Login,
	)
	s.Equal(http.StatusUnauthorized, resp.Code, "response body: %s", resp.Body.String())
}

func (s *IntegrationSuite) TestLogin_BlockedAfterTooManyFailedAttempts() {
	registerReq := models.RegisterRequest{Login: "throttled_user", Password: "1234567890123456"}
	regResp := s.PerformRequest(http.MethodPost, ServicePrefix+"/reg", registerReq, s.handler.Register)
	s.Require().Equal(http.StatusCreated, regResp.Code, "setup: registration must succeed")

	wrongLoginReq := models.LoginRequest{Login: registerReq.Login, Password: "wrongpassword12345"}

	maxAttempts := s.cfg.ThrottleCfg.MaxAttempts

	for i := range maxAttempts {
		resp := s.PerformRequest(http.MethodPost, ServicePrefix+"/login", wrongLoginReq, s.handler.Login)
		s.Equal(http.StatusUnauthorized, resp.Code, "attempt %d: response body: %s", i+1, resp.Body.String())
	}

	blockedResp := s.PerformRequest(http.MethodPost, ServicePrefix+"/login", wrongLoginReq, s.handler.Login)
	s.Equal(http.StatusTooManyRequests, blockedResp.Code, "response body: %s", blockedResp.Body.String())

	// Even the correct password is rejected while blocked.
	correctLoginReq := models.LoginRequest(registerReq)
	stillBlockedResp := s.PerformRequest(http.MethodPost, ServicePrefix+"/login", correctLoginReq, s.handler.Login)
	s.Equal(http.StatusTooManyRequests, stillBlockedResp.Code, "response body: %s", stillBlockedResp.Body.String())
}
