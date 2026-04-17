package dto

import "time"

// CreateBatchRequest defines the JSON payload for creating a batch settlement.
type CreateBatchRequest struct {
	SettlementDate string `json:"settlement_date" validate:"required"`
	MerchantID     string `json:"merchant_id" validate:"required"`
}

// BatchSettlementResponse defines the JSON response for a batch settlement.
type BatchSettlementResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	CreatedAtTime time.Time `json:"-"` // Internal use for sorting/filtering if needed
}
