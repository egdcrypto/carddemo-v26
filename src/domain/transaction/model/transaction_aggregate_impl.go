package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/transaction/command"
)

// Handle executes commands on the Transaction aggregate.
func (t *Transaction) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.SubmitTransactionCmd:
		return t.handleSubmit(c)
	case command.ReverseTransactionCmd:
		return t.handleReverse(c)
	default:
		return errors.New("unknown command")
	}
}

func (t *Transaction) handleSubmit(cmd command.SubmitTransactionCmd) error {
	// Domain Validation
	if cmd.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if cmd.AccountStatus != "Active" {
		return errors.New("account is not active")
	}

	// Apply state changes
	// (State is already set in NewTransaction, but we might update status/logic here)
	t.Status = "Submitted"
	return nil
}

func (t *Transaction) handleReverse(cmd command.ReverseTransactionCmd) error {
	// Domain Validation
	if t.Status == "Reversed" {
		return errors.New("transaction already reversed")
	}

	// Apply state changes
	t.Status = "Reversed"
	return nil
}
