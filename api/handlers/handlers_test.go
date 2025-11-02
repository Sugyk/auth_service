package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	checkJwtRoutePath := "/check_token"

	var req *http.Request
	var rr *httptest.ResponseRecorder
	handler := NewAPIHandler(nil, nil).CheckJWT()

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, checkJwtRoutePath, nil)

	test_login := "test_login"
	jwt, _ := CreateJWT(test_login)

	req.Header.Add("Authorization", "Bearer "+jwt)

	handler.ServeHTTP(rr, req)
	// Check with Authorization header
	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf(`Error of checking jwt
		Request Headers: %v
		Got: %d
		Expected: %d`, req.Header, rr.Result().StatusCode, http.StatusOK)
	}

	// Check without token in Authorization header
	req = httptest.NewRequest(http.MethodGet, checkJwtRoutePath, nil)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf(`Error of checking jwt
		Request Headers: %v
		Got: %d
		Expected: %d`, req.Header, rr.Result().StatusCode, http.StatusUnauthorized)
	}

	// Put wrong jwt to Authorization header
	req = httptest.NewRequest(http.MethodGet, checkJwtRoutePath, nil)
	req.Header.Set("Authorization", "Bearer "+"invalid_jwt")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf(`Error of checking jwt
		Request Headers: %v
		Got: %d
		Expected: %d`, req.Header, rr.Result().StatusCode, http.StatusUnauthorized)
	}

}
