package models

import "errors"

var (
	ErrDuplicate     = errors.New("duplicate")
	ErrLoginNotFound = errors.New("login not found")
)
