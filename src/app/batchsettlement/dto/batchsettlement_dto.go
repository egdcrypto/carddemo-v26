package dto

// CreateBatchSettlementRequest defines the JSON payload for creating a settlement.
type CreateBatchSettlementRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// BatchSettlementResponse defines the JSON response for a settlement.
type BatchSettlementResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
