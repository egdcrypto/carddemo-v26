package service

import "errors"

var (
	// ErrProfileNotFound is returned when a profile cannot be found.
	ErrProfileNotFound = errors.New("profile not found")
)
