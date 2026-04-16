package model

import (
	"time"
)

// Transaction represents the transaction aggregate.
type Transaction struct {
	ID              string
	AccountID       string
	CardID          string
	Amount          float64
	TransactionType string
	Status          string // e.g., "Submitted", "Reversed"
	CreatedAt       int64
	Version         int
}

// NewTransaction creates a new Transaction aggregate.
func NewTransaction(id, accountID, cardID string, amount float64, txType string) *Transaction {
	return &Transaction{
		ID:              id,
		AccountID:       accountID,
		CardID:          cardID,
		Amount:          amount,
		TransactionType: txType,
		Status:          "Submitted",
		CreatedAt:       time.Now().Unix(),
		Version:         1,
	}
}
