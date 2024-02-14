package postgres

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInternalServer = errors.New("internal server error")
	ErrLoginExists    = errors.New("login already exists")
)
