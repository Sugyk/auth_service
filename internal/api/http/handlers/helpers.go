package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

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
