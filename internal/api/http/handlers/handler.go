package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/pkg/logger"
)

type Service interface {
	Register(ctx context.Context, login string, password string) error
}

type Handler struct {
	service Service
	logger  logger.Logger
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.handleError(ctx, w, models.NewValidationErr("error decoding body"))
		return
	}
	if err := reqBody.Validate(); err != nil {
		h.handleError(ctx, w, err)
		return
	}

	serviceErr := h.service.Register(ctx, reqBody.Login, reqBody.Password)
	if serviceErr != nil {
		h.handleError(ctx, w, serviceErr)
		return
	}
	h.sendJSON(ctx, w, http.StatusCreated, nil)
}
