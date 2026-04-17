package dto

// UpdateUserProfileRequest represents the JSON payload for updating a user profile.
// The profile is owned by the account in this context.
type UpdateUserProfileRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

// UserProfileResponse represents the JSON response for a user profile.
type UserProfileResponse struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
