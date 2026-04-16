package dto

// CardPolicyResponse represents the JSON response for a card policy.
type CardPolicyResponse struct {
	ID              string `json:"id"`
	AccountID       string `json:"account_id"`
	DailyLimit      int    `json:"daily_limit"`
	SingleTxnLimit  int    `json:"single_txn_limit"`
	ActiveCountries []string `json:"active_countries"`
}

// UpdateCardPolicyRequest represents the JSON payload for updating a policy.
type UpdateCardPolicyRequest struct {
	DailyLimit     *int     `json:"daily_limit,omitempty" validate:"omitempty,min=0"`
	SingleTxnLimit *int     `json:"single_txn_limit,omitempty" validate:"omitempty,min=0"`
	ActiveCountries []string `json:"active_countries,omitempty" validate:"omitempty,dive,alpha"`
}
