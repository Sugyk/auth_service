package handlers

import (
	"context"
	"net/http"

	"github.com/Sugyk/auth_service/internal/models"
)

var codeToStatusMap = map[models.ErrorCode]int{
	models.CodeInternalError:   http.StatusInternalServerError,
	models.CodeErrDuplicate:    http.StatusConflict,
	models.CodeValidationError: http.StatusBadRequest,
}

func (h *Handler) mapErrToHttpStatus(ctx context.Context, errCode models.ErrorCode) int {
	status, ok := codeToStatusMap[errCode]
	if ok {
		return status
	}

	h.logger.Warn(ctx, "unknown error code encountered — falling back to 500",
		"code", errCode,
		"hint", "add this code to codeToStatus map in internal/api/http/helpers.go",
	)
	return http.StatusInternalServerError
}

func (h *Handler) handleError(ctx context.Context, w http.ResponseWriter, err error) {
	appErr, ok := models.AsAppError(err)
	if !ok {
		h.sendJSON(ctx, w, http.StatusInternalServerError, models.NewInternalErr())
		h.logger.Error(ctx, "internal error", "error", err.Error())
		return
	}

	httpStatus := h.mapErrToHttpStatus(ctx, appErr.ErrCode)

	h.sendJSON(ctx, w, httpStatus, appErr)
}
