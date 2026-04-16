package dto

// --- Request DTOs ---

// CreateAccountRequest represents the JSON payload for creating an account.
type CreateAccountRequest struct {
	UserProfileID string `json:"user_profile_id" validate:"required"`
	AccountType   string `json:"account_type" validate:"required,oneof=SAVINGS CHECKING CREDIT"`
	InitialStatus string `json:"initial_status,omitempty" validate:"omitempty,oneof=ACTIVE SUSPENDED"`
}

// UpdateAccountRequest represents the JSON payload for updating account details (e.g. status).
type UpdateAccountRequest struct {
	Status string `json:"status" validate:"required,oneof=ACTIVE SUSPENDED CLOSED"`
	Reason string `json:"reason" validate:"required,max=255"`
}

// --- Response DTOs ---

// AccountResponse is the standard JSON output for an account.
type AccountResponse struct {
	ID            string `json:"id"`
	UserProfileID string `json:"user_profile_id"`
	Status        string `json:"status"`
	AccountType   string `json:"account_type"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// ErrorResponse represents a structured error message returned to the client.
type ErrorResponse struct {
	Error string `json:"error"`
}
