package dto

// CreateAccountRequest represents the JSON payload for creating an account.
type CreateAccountRequest struct {
	UserProfileID string `json:"user_profile_id" validate:"required"`
	AccountType   string `json:"account_type" validate:"required"`
	Status        string `json:"status" validate:"required"`
}

// UpdateAccountStatusRequest represents the JSON payload for updating an account status.
type UpdateAccountStatusRequest struct {
	NewStatus string `json:"new_status" validate:"required"`
	Reason    string `json:"reason" validate:"required"`
}

// AccountResponse represents the JSON response for an account.
type AccountResponse struct {
	ID            string `json:"id"`
	UserProfileID string `json:"user_profile_id"`
	AccountType   string `json:"account_type"`
	Status        string `json:"status"`
	Version       int    `json:"version"`
}

// ErrorResponse represents an error message returned to the client.
type ErrorResponse struct {
	Error string `json:"error"`
}
