package model

import (
	"time"

	"github.com/google/uuid"
)

// Account represents the Account Aggregate.
type Account struct {
	ID            string
	UserProfileID string
	Status        string
	AccountType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Version       int
}

// NewAccount creates a new Account aggregate.
func NewAccount(userProfileID, accountType, status string) (*Account, error) {
	return &Account{
		ID:            uuid.New().String(),
		UserProfileID: userProfileID,
		AccountType:   accountType,
		Status:        status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Version:       0,
	}, nil
}
