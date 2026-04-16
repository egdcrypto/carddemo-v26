package dto

import "time"

// CreateSettlementRequest defines the JSON payload for creating a batch settlement.
type CreateSettlementRequest struct {
	MerchantID string `json:"merchant_id" validate:"required"`
	Currency   string `json:"currency" validate:"required,len=3"`
}

// SettlementResponse defines the JSON response for a settlement.
type SettlementResponse struct {
	ID         string    `json:"id"`
	MerchantID string    `json:"merchant_id"`
	Currency   string    `json:"currency"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
