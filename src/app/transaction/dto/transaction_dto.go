package dto

// CreateTransactionRequest defines the JSON payload for creating a transaction.
type CreateTransactionRequest struct {
	AccountID       string  `json:"account_id"`
	CardID          string  `json:"card_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

// VoidTransactionRequest defines the JSON payload for voiding a transaction.
type VoidTransactionRequest struct {
	Reason string `json:"reason"`
}

// TransactionResponse defines the JSON response for a transaction.
type TransactionResponse struct {
	ID              string  `json:"id"`
	AccountID       string  `json:"account_id"`
	CardID          string  `json:"card_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	Status          string  `json:"status"`
	WorkflowID      string  `json:"workflow_id,omitempty"`
}

// WorkflowResponse defines the response after triggering an async workflow.
type WorkflowResponse struct {
	WorkflowID string `json:"workflow_id"`
}
