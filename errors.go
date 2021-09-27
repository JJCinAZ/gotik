package gotik

import "errors"

var (
	ErrMissingId = errors.New("missing ID")
	ErrNotFound  = errors.New("not found")
)
