package dto

// CreateAccountRequest defines the JSON payload for creating an account.
type CreateAccountRequest struct {
	UserProfileID string `json:"user_profile_id" validate:"required"`
	Status        string `json:"status" validate:"required,oneof=ACTIVE SUSPENDED CLOSED"`
	AccountType   string `json:"account_type" validate:"required,oneof=CHECKING SAVINGS"`
}

// UpdateAccountRequest defines the JSON payload for updating an account.
type UpdateAccountRequest struct {
	Status string `json:"status" validate:"required,oneof=ACTIVE SUSPENDED CLOSED"`
	Reason string `json:"reason" validate:"required"`
}

// UpdateUserProfileRequest defines the JSON payload for linking a user profile.
type UpdateUserProfileRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

// AccountResponse defines the JSON response for an account.
type AccountResponse struct {
	ID            string `json:"id"`
	UserProfileID string `json:"user_profile_id"`
	Status        string `json:"status"`
	AccountType   string `json:"account_type"`
	Version       int    `json:"version"`
}

// UserProfileResponse defines the JSON response for a user profile.
type UserProfileResponse struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// ErrorResponse defines the standard error JSON response.
type ErrorResponse struct {
	Error string `json:"error"`
}
