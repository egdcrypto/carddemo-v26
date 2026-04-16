package dto

// IssueCardRequest represents the JSON payload for issuing a new card.
type IssueCardRequest struct {
	AccountID string `json:"account_id" validate:"required"`
	CardType  string `json:"card_type" validate:"required,oneof=Virtual Physical"`
}

// CardResponse represents the JSON response for a card.
type CardResponse struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Status    string `json:"status"`
	CardType  string `json:"card_type"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
}

// UpdateCardStatusRequest represents the JSON payload for updating card status.
type UpdateCardStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=Active Blocked Closed"`
}

// ActivateCardRequest represents the JSON payload for activating a card.
type ActivateCardRequest struct {
	PIN string `json:"pin" validate:"required,len=4"`
}
