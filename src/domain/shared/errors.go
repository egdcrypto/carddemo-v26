package shared

import "errors"

var (
	// ErrNotFound is returned when an entity is not found.
	ErrNotFound = errors.New("entity not found")
)
