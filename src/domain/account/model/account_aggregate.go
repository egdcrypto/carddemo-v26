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
	// Placeholder implementation for Red phase.
	// Logic will be implemented in the Green phase.
	return nil, nil
}
