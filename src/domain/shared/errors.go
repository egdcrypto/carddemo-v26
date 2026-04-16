package shared

import "errors"

var (
	// ErrUnknownCommand is returned when an unregistered command is executed.
	ErrUnknownCommand = errors.New("unknown command")

	// ErrInvalidAccountStatus is returned when an account is not in a valid state for operations.
	ErrInvalidAccountStatus = errors.New("account status must be 'Pending' or 'Active' to process financial transactions")

	// ErrAccountClosed is returned when trying to modify a closed account.
	ErrAccountClosed = errors.New("account closure is irreversible and requires a zero balance")
)
