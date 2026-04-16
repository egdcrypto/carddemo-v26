package dto

// CreateSettlementRequest defines the JSON payload for creating a settlement batch.
type CreateSettlementRequest struct {
	MerchantID string  `json:"merchant_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

// SettlementResponse defines the JSON response for a settlement batch.
type SettlementResponse struct {
	ID         string  `json:"id"`
	MerchantID string  `json:"merchant_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}
