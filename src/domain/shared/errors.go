package shared

import "errors"

// Common domain errors.
var (
	// ErrUnknownCommand is returned when an aggregate receives a command it cannot handle.
	ErrUnknownCommand = errors.New("unknown command")

	// ErrValidation is returned when business rule validation fails.
	ErrValidation = errors.New("validation failed")
)
