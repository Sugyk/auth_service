package grpc_api

import (
	"github.com/Sugyk/auth_service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var codeToGRPCStatusMap = map[models.ErrorCode]codes.Code{
	models.CodeInternalError:    codes.Internal,
	models.CodeErrDuplicate:     codes.AlreadyExists,
	models.CodeValidationError:  codes.InvalidArgument,
	models.CodeWrongCredentials: codes.Unauthenticated,
}

func appErrorFrom(err error) *models.AppError {
	if appErr, ok := models.AsAppError(err); ok {
		return appErr
	}
	return models.NewInternalErr(err.Error())
}

func toGRPCError(appErr *models.AppError) error {
	code, ok := codeToGRPCStatusMap[appErr.ErrCode]
	if !ok {
		code = codes.Internal
	}

	return status.Error(code, appErr.Details)
}
