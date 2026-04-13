package models

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	CodeErrDuplicate    ErrorCode = "Duplicate login"
	CodeValidationError ErrorCode = "Validation error"
	CodeInternalError   ErrorCode = "Internal error"
)

// Error type that contains ready info about errors for user
// cause - parameter to use in logs for errors that user should not see
type AppError struct {
	ErrCode ErrorCode `json:"error"`
	Details string    `json:"details"`
}

func (a *AppError) Error() string {
	return fmt.Sprintf("%s: %s", a.ErrCode, a.Details)
}

func New(errCode ErrorCode, details string) *AppError {
	return &AppError{
		ErrCode: errCode,
		Details: details,
	}
}

func NewDuplicateLoginErr(login string) *AppError {
	return New(
		CodeErrDuplicate,
		fmt.Sprintf("Login %s is already exists", login),
	)
}

func NewValidationErr(message string) *AppError {
	return New(
		CodeValidationError,
		message,
	)
}

func NewInternalErr() *AppError {
	return New(CodeInternalError, "internal error")
}

// Helpers

// Helper function to check is error is already AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if ok := errors.As(err, &appErr); ok {
		return appErr, ok
	}
	return nil, false
}
