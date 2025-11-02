package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckJWT(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := NewAPIHandler(nil, nil).CheckJWT()
	req := httptest.NewRequest(http.MethodGet, "/check_token", nil)
	handler.ServeHTTP(rr, req)

	jwtSecret = []byte("secret")

	// if rr.Result().StatusCode !=
}
