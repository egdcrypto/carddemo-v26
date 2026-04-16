package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/event"
)

var (
	// ErrTransactionNotFound is returned when a transaction is not found.
	ErrTransactionNotFound = errors.New("transaction not found")
	// ErrInvalidAmount is returned when the amount is invalid.
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrInvalidStatus is returned when the status is invalid for the operation.
	ErrInvalidStatus = errors.New("invalid status")
)

// Transaction represents the Transaction Aggregate.
type Transaction struct {
	shared.AggregateBase
	AccountID       string
	CardID          string
	Amount          float64
	Type            string // debit, credit
	Status          string // pending, cleared, reversed
}

// NewTransaction creates a new Transaction aggregate.
func NewTransaction(id string) *Transaction {
	return &Transaction{
		AggregateBase: shared.AggregateBase{ID: id},
		Status:        "pending",
	}
}

// Execute handles commands.
func (t *Transaction) Execute(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.SubmitTransactionCmd:
		return t.handleSubmit(c)
	case command.ReverseTransactionCmd:
		return t.handleReverse(c)
	default:
		return errors.New("unknown command")
	}
}

func (t *Transaction) handleSubmit(c command.SubmitTransactionCmd) error {
	// In a real scenario, validation logic would be more robust.
	if c.Amount <= 0 {
		return ErrInvalidAmount
	}
	if c.AccountStatus != "Active" {
		return ErrInvalidStatus
	}

	t.AccountID = c.AccountID
	t.CardID = c.CardID
	t.Amount = c.Amount
	t.Type = c.TransactionType
	t.Status = "cleared" // Simulate immediate clearance

	e := event.TransactionSubmitted{
		TransactionID:   t.ID,
		AccountID:       t.AccountID,
		CardID:          t.CardID,
		Amount:          t.Amount,
		TransactionType: t.Type,
	}
	t.AddEvent(e)

	return nil
}

func (t *Transaction) handleReverse(c command.ReverseTransactionCmd) error {
	if t.Status == "reversed" {
		return errors.New("transaction already reversed")
	}
	if t.Status != "cleared" {
		return errors.New("cannot reverse uncleared transaction")
	}

	t.Status = "reversed"

	e := event.TransactionReversed{
		TransactionID: t.ID,
		Reason:        c.Reason,
	}
	t.AddEvent(e)

	return nil
}
