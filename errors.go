package gotik

import "errors"

var (
	ErrMissingId     = errors.New("missing ID")
	ErrNotFound      = errors.New("not found")
	ErrMissingChain  = errors.New("missing chain")
	ErrVersionTooOld = errors.New("RouterOS version too old")
)
