package command

// UpdateProfileCommand represents a request to update user details.
type UpdateProfileCommand struct {
	FirstName string
	LastName  string
	Email     string
}
