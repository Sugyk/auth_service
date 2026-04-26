package models

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	CodeErrDuplicate     ErrorCode = "Duplicate login"
	CodeValidationError  ErrorCode = "Validation error"
	CodeInternalError    ErrorCode = "Internal error"
	CodeWrongCredentials ErrorCode = "Wrong credentials"
)

// Error type that contains ready info about errors for user
// cause - parameter to use in logs for errors that user should not see
type AppError struct {
	ErrCode ErrorCode `json:"error"`
	Details string    `json:"details"`
	cause   string    `json:"-"`
}

func (a *AppError) Error() string {
	return fmt.Sprintf("%s: %s", a.ErrCode, a.Details)
}

func (e *AppError) Cause() string {
	return e.cause
}

func New(errCode ErrorCode, details string, cause string) *AppError {
	return &AppError{
		ErrCode: errCode,
		Details: details,
		cause:   cause,
	}
}

func NewDuplicateLoginErr(login string) *AppError {
	return New(
		CodeErrDuplicate,
		fmt.Sprintf("Login %s is already exists", login),
		"",
	)
}

func NewValidationErr(message string) *AppError {
	return New(
		CodeValidationError,
		message,
		"",
	)
}

func NewInternalErr(cause string) *AppError {
	return New(CodeInternalError, "internal error", cause)
}

func NewLoginNotFound(login string) *AppError {
	return New(
		CodeWrongCredentials,
		"There is incorrect login or password",
		fmt.Sprintf("login '%s' not found", login),
	)
}

func NewWrongPassword(login string) *AppError {
	return New(
		CodeWrongCredentials,
		"There is incorrect login or password",
		fmt.Sprintf("incorrect password for login '%s'", login),
	)
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
