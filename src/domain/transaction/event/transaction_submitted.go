package event

import "time"

// TransactionSubmitted is emitted when a transaction is successfully validated and applied.
type TransactionSubmitted struct {
	OccurredAt   time.Time
	AggregateID  string
	AccountID    string
	CardID       string
	Amount       float64
	Type         string
}

// NewTransactionSubmitted creates a new TransactionSubmitted event.
func NewTransactionSubmitted(aggregateID, accountID, cardID string, amount float64, t string) *TransactionSubmitted {
	return &TransactionSubmitted{
		OccurredAt:  time.Now(),
		AggregateID: aggregateID,
		AccountID:   accountID,
		CardID:      cardID,
		Amount:      amount,
		Type:        t,
	}
}
