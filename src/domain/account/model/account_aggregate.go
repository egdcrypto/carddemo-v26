package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	// ErrInvalidCommand is returned when an invalid command is applied.
	ErrInvalidCommand = errors.New("invalid command")
)

// Account represents the Account Aggregate.
type Account struct {
	shared.AggregateRoot
	ID            string
	UserProfileID string
	AccountType   string
	Status        string
	Version       int
}

// NewAccount creates a new Account aggregate.
func NewAccount(id string) *Account {
	return &Account{
		AggregateRoot: shared.AggregateRoot{},
		ID:            id,
		Version:       0,
	}
}

// Handle processes commands.
func (a *Account) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.OpenAccountCmd:
		return a.openAccount(c)
	case command.UpdateAccountStatusCmd:
		return a.updateStatus(c)
	default:
		return ErrInvalidCommand
	}
}

// openAccount handles account creation.
func (a *Account) openAccount(cmd command.OpenAccountCmd) error {
	// In a real scenario, check invariants (e.g. valid status, type).
	a.UserProfileID = cmd.UserProfileID
	a.AccountType = cmd.AccountType
	a.Status = cmd.InitialStatus

	// Create Event
	e := event.NewAccountOpened(a.ID, cmd)
	e.Payload.AccountID = a.ID
	e.Payload.UserProfileID = cmd.UserProfileID
	e.Payload.Status = cmd.InitialStatus
	e.Payload.AccountType = cmd.AccountType

	a.AddDomainEvent(e)
	return nil
}

// updateStatus handles status updates.
func (a *Account) updateStatus(cmd command.UpdateAccountStatusCmd) error {
	oldStatus := a.Status
	a.Status = cmd.NewStatus

	// Create Event
	e := event.NewAccountStatusUpdated(a.ID)
	e.Payload.AccountID = a.ID
	e.Payload.OldStatus = oldStatus
	e.Payload.NewStatus = cmd.NewStatus
	e.Payload.Reason = cmd.Reason

	a.AddDomainEvent(e)
	return nil
}
