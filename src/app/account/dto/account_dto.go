package dto

// Validation errors
var (
	ErrInvalidJSON        = "invalid JSON"
	ErrInvalidID          = "invalid account ID"
	ErrInvalidStatus      = "invalid status"
	ErrMissingName        = "name is required"
	ErrMissingEmail       = "email is required"
)

// AccountResponse represents the JSON response for an account.
type AccountResponse struct {
	ID            string `json:"id"`
	UserProfileID string `json:"user_profile_id"`
	Status        string `json:"status"`
	AccountType   string `json:"account_type"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CreateAccountRequest represents the JSON request to create an account.
type CreateAccountRequest struct {
	UserProfileID string `json:"user_profile_id"`
	Status        string `json:"status"` // Should validate against domain logic if needed
	AccountType   string `json:"account_type"`
}

// Validate checks the request fields.
func (r *CreateAccountRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.UserProfileID == "" {
		errs["user_profile_id"] = ErrMissingName // Reusing generic error or define specific
	}
	if r.AccountType == "" {
		errs["account_type"] = ErrMissingName
	}
	// Add more validation as per AC
	return errs
}

// UpdateAccountRequest represents the JSON request to update an account.
type UpdateAccountRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// Validate checks the request fields.
func (r *UpdateAccountRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Status == "" {
		errs["status"] = ErrMissingName
	}
	return errs
}

// UpdateUserProfileRequest represents the JSON request to update a user profile.
type UpdateUserProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Validate checks the request fields.
func (r *UpdateUserProfileRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.FirstName == "" {
		errs["first_name"] = ErrMissingName
	}
	if r.LastName == "" {
		errs["last_name"] = ErrMissingName
	}
	if r.Email == "" {
		errs["email"] = ErrMissingEmail
	}
	return errs
}

// UserProfileResponse represents the JSON response for a user profile.
type UserProfileResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
