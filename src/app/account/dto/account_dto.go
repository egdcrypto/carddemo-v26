package dto

// "Acceptance Criteria: Request and response types are documented and versioned"

// CreateAccountRequest defines the JSON payload for creating an account.
type CreateAccountRequest struct {
	UserProfileID string `json:"user_profile_id" validate:"required"`
	Status        string `json:"status" validate:"required,oneof=Active Pending Suspended"`
	AccountType   string `json:"account_type" validate:"required,oneof=Checking Savings Credit"`
}

// UpdateAccountStatusRequest defines the JSON payload for updating status.
type UpdateAccountStatusRequest struct {
	NewStatus string `json:"new_status" validate:"required,oneof=Active Pending Suspended Closed"`
	Reason    string `json:"reason" validate:"required,max=200"`
}

// AccountResponse defines the JSON response for an account.
type AccountResponse struct {
	ID            string `json:"id"`
	UserProfileID string `json:"user_profile_id"`
	Status        string `json:"status"`
	AccountType   string `json:"account_type"`
	Version       int    `json:"version"`
}

// ErrorResponse defines the standard error JSON structure.
type ErrorResponse struct {
	Error string `json:"error"`
}