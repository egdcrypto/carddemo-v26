package model

import (
	"time"
)

// BatchSettlement represents the settlement aggregate.
type BatchSettlement struct {
	ID        string
	MerchantID string
	Currency   string
	Status     string // e.g., "Pending", "Completed"
	CreatedAt time.Time
	Version   int
}

// NewBatchSettlement creates a new BatchSettlement aggregate.
func NewBatchSettlement(id, merchantID, currency string) *BatchSettlement {
	return &BatchSettlement{
		ID:        id,
		MerchantID: merchantID,
		Currency:  currency,
		Status:    "Pending",
		CreatedAt: time.Now(),
		Version:   1,
	}
}
