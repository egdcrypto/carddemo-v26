package dto

// LinkUserToAccountRequest defines the JSON payload for linking a profile.
// Note: This is likely handled via the UserProfile aggregate logic,
// but the DTO lives here for the REST adapter.
type LinkUserToAccountRequest struct {
	AccountID string `json:"account_id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// UserProfileResponse defines the JSON response for a user profile.
type UserProfileResponse struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
