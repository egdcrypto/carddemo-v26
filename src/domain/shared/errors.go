package rest

import "errors"

// Common errors used by handlers.
var (
	// ErrNotFound is returned when a requested resource is not found.
	ErrNotFound = errors.New("resource not found")
)
