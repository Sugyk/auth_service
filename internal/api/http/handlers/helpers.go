package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

func (h *Handler) sendJSON(ctx context.Context, w http.ResponseWriter, httpStatus int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.logger.Error(ctx, "error send JSON", "error", err)
	}
}
