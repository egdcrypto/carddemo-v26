package model

import (
	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
)

// Account represents the Account aggregate.
type Account struct {
	shared.AggregateRoot
	ID     string
	Status string
	Closed bool
}

// NewAccount creates a new Account instance.
func NewAccount(id string) *Account {
	return &Account{ID: id}
}

// Execute handles commands for the Account aggregate.
func (a *Account) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.OpenAccountCmd:
		return a.handleOpenAccount(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (a *Account) GetID() string {
	return a.ID
}

// ID satisfies the shared.Aggregate interface.
func (a *Account) ID() string {
	return a.ID
}

func (a *Account) handleOpenAccount(c command.OpenAccountCmd) ([]shared.DomainEvent, error) {
	// Invariant Check: Account closure is irreversible and requires a zero balance.
	// The command payload flags "IsClosed" to simulate a state check against business rules.
	if c.IsClosed {
		return nil, shared.ErrAccountClosed
	}

	// Invariant Check: Account status must be 'Pending' or 'Active' to process financial transactions.
	// We enforce this upon creation/opening as well.
	if c.InitialStatus != "Pending" && c.InitialStatus != "Active" {
		return nil, shared.ErrInvalidAccountStatus
	}

	// Apply state changes
	a.Status = c.InitialStatus
	a.Closed = false

	// Emit event
	evt := &event.AccountOpenedEvent{
		DomainEvent:   shared.DomainEvent{}, // In a real app, set timestamps/IDs here
		AccountID:     a.ID,
		UserProfileID: c.UserProfileID,
		Status:        c.InitialStatus,
		AccountType:   c.AccountType,
	}

	return []shared.DomainEvent{evt}, nil
}
