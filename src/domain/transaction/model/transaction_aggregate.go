package model

import (
	"errors"
	"fmt"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/event"
)

// Transaction represents the Transaction aggregate.
type Transaction struct {
	shared.AggregateRoot
	ID string
}

// NewTransaction creates a new Transaction instance.
func NewTransaction(id string) *Transaction {
	return &Transaction{ID: id}
}

// Execute handles commands for the Transaction aggregate.
func (t *Transaction) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.SubmitTransactionCmd:
		return t.handleSubmitTransaction(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleSubmitTransaction validates and applies the transaction submission.
func (t *Transaction) handleSubmitTransaction(cmd command.SubmitTransactionCmd) ([]shared.DomainEvent, error) {
	// Rule 1: Transaction amount must be strictly greater than zero
	if cmd.Amount <= 0 {
		return nil, errors.New("transaction amount must be strictly greater than zero")
	}

	// Rule 2: Account must be in 'Active' status to accept debit or credit transactions
	if cmd.AccountStatus != "Active" {
		return nil, fmt.Errorf("account must be in 'Active' status to accept debit or credit transactions, current status: %s", cmd.AccountStatus)
	}

	// Create the domain event
	evt := event.NewTransactionSubmitted(t.ID, cmd.AccountID, cmd.CardID, cmd.Amount, cmd.Type)

	// Return event slice
	return []shared.DomainEvent{evt}, nil
}

// ID returns the aggregate ID.
func (t *Transaction) GetID() string {
	return t.ID
}
