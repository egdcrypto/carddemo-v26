package dto

// IssueCardRequest defines the JSON payload for issuing a new card.
type IssueCardRequest struct {
	CardType       string            `json:"card_type" validate:"required,oneof=virtual physical"`
	SpendingLimits map[string]int    `json:"spending_limits" validate:"dive,keys,required,endkeys,required,min=0"`
	Meta           map[string]string `json:"meta,omitempty"`
}

// UpdateCardStatusRequest defines the JSON payload for updating a card's status.
type UpdateCardStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active blocked suspended closed lost_stolen"`
	Reason string `json:"reason,omitempty" validate:"max=255"`
}

// ActivateCardRequest defines the JSON payload for activating a card.
type ActivateCardRequest struct {
	ActivationCode string `json:"activation_code" validate:"required,len=6"`
}

// UpdateCardPolicyRequest defines the JSON payload for updating policy limits.
type UpdateCardPolicyRequest struct {
	DailyLimit int `json:"daily_limit" validate:"required,min=0"`
	WeeklyLimit int `json:"weekly_limit" validate:"min=0"`
}

// CardResponse represents the HTTP response for a card.
type CardResponse struct {
	ID             string            `json:"id"`
	AccountID      string            `json:"account_id"`
	CardType       string            `json:"card_type"`
	Status         string            `json:"status"`
	SpendingLimits map[string]int    `json:"spending_limits"`
	MaskedPAN      string            `json:"masked_pan"`
	CreatedAt      string            `json:"created_at"`
}

// CardPolicyResponse represents the HTTP response for a card policy.
type CardPolicyResponse struct {
	ID          string `json:"id"`
	CardID      string `json:"card_id"`
	DailyLimit  int    `json:"daily_limit"`
	WeeklyLimit int    `json:"weekly_limit"`
	IsActive    bool   `json:"is_active"`
}

// ErrorResponse defines the standard error JSON structure.
type ErrorResponse struct {
	Error string `json:"error"`
}
