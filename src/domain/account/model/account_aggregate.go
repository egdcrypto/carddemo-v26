package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/google/uuid"
)

var (
	// ErrInvalidCommand is returned when the command is not applicable.
	ErrInvalidCommand = errors.New("invalid command")
)

// Account is the aggregate root.
type Account struct {
	shared.AggregateRoot
	ID            string
	UserProfileID string
	Status        string
	AccountType   string
	Version       int
}

// NewAccount creates a new Account aggregate.
func NewAccount(profileID, status, accountType string) *Account {
	return &Account{
		ID:            uuid.New().String(),
		UserProfileID: profileID,
		Status:        status,
		AccountType:   accountType,
		Version:       1,
		AggregateRoot: shared.AggregateRoot{},
	}
}

// Execute handles commands.
func (a *Account) Execute(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.OpenAccountCmd:
		return a.handleOpenAccount(c)
	case command.UpdateAccountStatusCmd:
		return a.handleUpdateStatus(c)
	default:
		return ErrInvalidCommand
	}
}

func (a *Account) handleOpenAccount(cmd command.OpenAccountCmd) error {
	// Ideally, logic to check initial state goes here.
	// Since this is a constructor flow in the handler for simplicity, we just record the event.
	e := event.NewAccountOpened(a.ID, cmd)
	// Populate event payload
	e.Payload.AccountID = a.ID
	e.Payload.UserProfileID = a.UserProfileID
	e.Payload.Status = a.Status
	e.Payload.AccountType = a.AccountType

	a.AddEvent(e)
	return nil
}

func (a *Account) handleUpdateStatus(cmd command.UpdateAccountStatusCmd) error {
	oldStatus := a.Status
	// Simple business logic validation
	if a.Status == cmd.NewStatus {
		return nil // No-op
	}

	a.Status = cmd.NewStatus

	e := event.NewAccountStatusUpdated(a.ID)
	e.Payload.AccountID = a.ID
	e.Payload.OldStatus = oldStatus
	e.Payload.NewStatus = cmd.NewStatus
	e.Payload.Reason = cmd.Reason

	a.AddEvent(e)
	return nil
}
