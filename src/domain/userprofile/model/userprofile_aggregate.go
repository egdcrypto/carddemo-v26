package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/userprofile/command"
	"github.com/carddemo/project/src/domain/userprofile/event"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	ErrInvalidCommand = errors.New("invalid command")
)

// UserProfile represents the UserProfile Aggregate.
type UserProfile struct {
	shared.AggregateRoot
	ID        string
	AccountID string
	FirstName string
	LastName  string
	Email     string
	Version   int
}

// NewUserProfile creates a new UserProfile aggregate.
func NewUserProfile(id string, accountID string) *UserProfile {
	return &UserProfile{
		AggregateRoot: shared.AggregateRoot{},
		ID:            id,
		AccountID:     accountID,
		Version:       0,
	}
}

// Handle processes commands.
func (u *UserProfile) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.UpdateProfileCommand:
		return u.updateProfile(c)
	default:
		return ErrInvalidCommand
	}
}

// updateProfile handles profile updates.
func (u *UserProfile) updateProfile(cmd command.UpdateProfileCommand) error {
	u.FirstName = cmd.FirstName
	u.LastName = cmd.LastName
	u.Email = cmd.Email

	// Create Event
	e := event.NewUserProfileUpdated(u.ID)
	e.Payload.UserProfileID = u.ID
	e.Payload.AccountID = u.AccountID
	e.Payload.FirstName = cmd.FirstName
	e.Payload.LastName = cmd.LastName
	e.Payload.Email = cmd.Email

	u.AddDomainEvent(e)
	return nil
}
