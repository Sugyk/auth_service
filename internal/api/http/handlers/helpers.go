package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Sugyk/auth_service/internal/models"
)

type Request interface {
	Validate() error
}

func (h *Handler) readRequestBody(r *http.Request, d Request) error {
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return models.NewValidationErr("error decoding body")
	}
	if err := d.Validate(); err != nil {
		return err
	}
	return nil
}

func (h *Handler) sendJSON(ctx context.Context, w http.ResponseWriter, httpStatus int, body any) {
	data, err := json.Marshal(body)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(data)
}
