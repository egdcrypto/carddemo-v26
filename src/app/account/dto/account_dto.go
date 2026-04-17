package dto

// CreateAccountRequest defines the JSON payload for creating an account.
type CreateAccountRequest struct {
	UserProfileID string `json:"user_profile_id" validate:"required"`
	AccountType   string `json:"account_type" validate:"required"`
	Status        string `json:"status" validate:"required"`
}

// UpdateAccountStatusRequest defines the JSON payload for updating status.
type UpdateAccountStatusRequest struct {
	NewStatus string `json:"new_status" validate:"required"`
	Reason    string `json:"reason" validate:"required"`
}

// AccountResponse defines the JSON response for an account.
type AccountResponse struct {
	ID            string `json:"id"`
	UserProfileID string `json:"user_profile_id"`
	Status        string `json:"status"`
	AccountType   string `json:"account_type"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Version       int    `json:"version"`
}
