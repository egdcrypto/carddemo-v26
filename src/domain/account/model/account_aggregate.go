package model

import (
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	// ErrInvalidStatus is returned when an invalid status transition is attempted.
	ErrInvalidStatus = errors.New("invalid status")
	// ErrEmptyReason is returned if a reason is not provided for status change.
	ErrEmptyReason = errors.New("reason cannot be empty")
)

// Account is the aggregate root for the Account domain.
type Account struct {
	shared.AggregateRoot
	UserProfileID string
	Status        string
	AccountType   string
}

// NewAccount creates a new Account aggregate.
func NewAccount(cmd command.OpenAccountCmd) *Account {
	id := shared.GenerateUUID() // Assuming helper exists in shared
	a := &Account{
		AggregateRoot: shared.AggregateRoot{
			ID:      id,
			Version: 1,
		},
		UserProfileID: cmd.UserProfileID,
		Status:        cmd.InitialStatus,
		AccountType:   cmd.AccountType,
	}

	// Record Domain Event
	e := event.NewAccountOpened(id, cmd)
	e.Payload.AccountID = id
	e.Payload.UserProfileID = cmd.UserProfileID
	e.Payload.Status = cmd.InitialStatus
	e.Payload.AccountType = cmd.AccountType
	a.AddEvent(e)

	return a
}

// Handle processes commands for the Account aggregate.
func (a *Account) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.UpdateAccountStatusCmd:
		return a.updateStatus(c)
	default:
		return errors.New("unknown command")
	}
}

func (a *Account) updateStatus(cmd command.UpdateAccountStatusCmd) error {
	// Business Logic Validation
	if cmd.NewStatus == a.Status {
		return nil // No-op
	}

	// Check valid transitions (simplified)
	if cmd.NewStatus == "" {
		return ErrInvalidStatus
	}

	a.Status = cmd.NewStatus

	// Record Domain Event
	e := event.NewAccountStatusUpdated(a.ID)
	e.Payload.AccountID = a.ID
	e.Payload.OldStatus = a.Status // Simplification: Should store old before change
	e.Payload.NewStatus = cmd.NewStatus
	e.Payload.Reason = cmd.Reason
	a.AddEvent(e)

	a.Version++
	return nil
}
