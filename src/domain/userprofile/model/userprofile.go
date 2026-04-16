package model

import (
	"time"

	"github.com/google/uuid"
)

// UserProfile represents the UserProfile Aggregate.
type UserProfile struct {
	ID        string
	AccountID string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

// NewUserProfile creates a new UserProfile aggregate.
func NewUserProfile(firstName, lastName, email string) (*UserProfile, error) {
	return &UserProfile{
		ID:        uuid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   0,
	}, nil
}
