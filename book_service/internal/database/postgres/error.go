package postgres

import "errors"

var (
	ErrAuthorExists   = errors.New("author exists")
	ErrInternalServer = errors.New("internal server error")
	ErrAuthorNotFound = errors.New("author not found")
	ErrBookExists     = errors.New("book exists")
	ErrBookNotFound   = errors.New("book not found")
)
