package domain

import "errors"

var (
	ErrValidation   = errors.New("validation")
	ErrInternal     = errors.New("internal")
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotAllowed   = errors.New("not allowed")
)
