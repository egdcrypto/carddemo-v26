package dto

// --- Request DTOs ---

// UpdateUserProfileRequest represents the JSON payload for updating a user profile.
type UpdateUserProfileRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
}

// --- Response DTOs ---

// UserProfileResponse is the standard JSON output for a user profile.
type UserProfileResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	AccountID string `json:"account_id"`
}
