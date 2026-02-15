package gotik

import "errors"

var (
	ErrMissingId          = errors.New("missing ID")
	ErrNotFound           = errors.New("not found")
	ErrMissingChain       = errors.New("missing chain")
	ErrVersionTooOld      = errors.New("RouterOS version too old")
	ErrIPorPrefixNotValid = errors.New("IP must be a valid IPv4/IPv6 address or subnet")
)
