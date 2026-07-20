//go:generate mockgen -destination=service_mock.go -source=handler.go -package=handlers

package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/pkg/logger"
)

type Service interface {
	Register(ctx context.Context, login string, password string) error
	Login(ctx context.Context, login string, password string) (string, error)
}

type Handler struct {
	service Service
	logger  logger.Logger
}

func NewHandler(service Service, logger logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a user with the given login and password. Password must be 16+ characters.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.RegisterRequest  true  "Registration payload"
// @Success      201      {object}  models.RegisterResponse201
// @Failure      400      {object}  models.AppError  "validation error"
// @Failure      409      {object}  models.AppError  "login already exists"
// @Failure      500      {object}  models.AppError  "internal error"
// @Router       /auth/reg [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody models.RegisterRequest
	if err := h.readRequestBody(r, &reqBody); err != nil {
		h.handleError(ctx, w, err)
		return
	}

	serviceErr := h.service.Register(ctx, reqBody.Login, reqBody.Password)
	if serviceErr != nil {
		h.handleError(ctx, w, serviceErr)
		return
	}

	resp := models.RegisterResponse201{
		Message: fmt.Sprintf("user with login '%s' created", reqBody.Login),
	}

	h.sendJSON(ctx, w, http.StatusCreated, resp)
}

// Login godoc
// @Summary      Log in
// @Description  Authenticates a user and returns a signed JWT access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.LoginRequest  true  "Login payload"
// @Success      200      {object}  models.LoginResponse200
// @Failure      400      {object}  models.AppError  "validation error"
// @Failure      401      {object}  models.AppError  "incorrect login or password"
// @Failure      429      {object}  models.AppError  "too many failed attempts"
// @Failure      500      {object}  models.AppError  "internal error"
// @Router       /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody models.LoginRequest
	if err := h.readRequestBody(r, &reqBody); err != nil {
		h.handleError(ctx, w, err)
		return
	}

	accessToken, serviceErr := h.service.Login(ctx, reqBody.Login, reqBody.Password)

	if serviceErr != nil {
		h.handleError(ctx, w, serviceErr)
		return
	}

	resp := models.LoginResponse200{
		AccessToken: accessToken,
	}

	h.sendJSON(ctx, w, http.StatusOK, resp)
}
