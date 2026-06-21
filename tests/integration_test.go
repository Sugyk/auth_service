package integration_tests

import (
	"net/http"

	"github.com/Sugyk/auth_service/internal/models"
)

const ServicePrefix = "/api/v1/auth"

func (s *IntegrationSuite) TestRegister() {
	resp := s.PerformRequest(
		http.MethodPost,
		ServicePrefix+"/register",
		models.RegisterRequest{
			Login:    "testuser",
			Password: "StrongPass12345678!",
		},
		s.handler.Register,
	)

	s.Equal(http.StatusCreated, resp.Code)
}
