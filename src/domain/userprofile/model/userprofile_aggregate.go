package model

import (
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/userprofile/command"
	"github.com/carddemo/project/src/domain/userprofile/event"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
)

// UserProfile is the aggregate root.
type UserProfile struct {
	shared.AggregateRoot
	ID        string
	AccountID string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUserProfile creates a new UserProfile aggregate.
func NewUserProfile(id, accountID, firstName, lastName, email string) *UserProfile {
	return &UserProfile{
		AggregateRoot: shared.AggregateRoot{Version: 0},
		ID:            id,
		AccountID:     accountID,
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// Execute handles commands.
func (u *UserProfile) Execute(cmd interface{}) error {
	switch c := cmd.(type) {
	case *command.LinkUserToAccountCommand:
		return u.linkUser(c)
	case *command.RegisterUserCommand:
		return u.registerUser(c)
	default:
		return shared.ErrUnknownCommand
	}
}

func (u *UserProfile) linkUser(cmd *command.LinkUserToAccountCommand) error {
	// If it's a new profile, set basic info. If existing, update.
	// Assuming this handler is used for the PUT /profile endpoint which effectively upserts/links.
	u.FirstName = cmd.FirstName
	u.LastName = cmd.LastName
	u.UpdatedAt = time.Now()

	evt := event.NewUserProfileLinked(u.ID)
	evt.Payload.UserProfileID = u.ID
	evt.Payload.AccountID = u.AccountID
	evt.Payload.FirstName = u.FirstName
	evt.Payload.LastName = u.LastName

	u.AddEvent(evt)
	return nil
}

func (u *UserProfile) registerUser(cmd *command.RegisterUserCommand) error {
	u.FirstName = cmd.FirstName
	u.LastName = cmd.LastName
	u.Email = cmd.Email
	u.UpdatedAt = time.Now()

	evt := event.NewUserRegistered(u.ID)
	evt.Payload.UserProfileID = u.ID
	evt.Payload.Email = u.Email

	u.AddEvent(evt)
	return nil
}
