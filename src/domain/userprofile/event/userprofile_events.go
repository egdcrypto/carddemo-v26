package event

import (
	"github.com/google/uuid"
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// UserProfileUpdated is published when user details change.
type UserProfileUpdated struct {
	shared.CloudEventEnvelope
	Payload struct {
		UserProfileID string `json:"user_profile_id"`
		AccountID     string `json:"account_id"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Email         string `json:"email"`
	} `json:"data"`
}

// NewUserProfileUpdated creates a new UserProfileUpdated event.
func NewUserProfileUpdated(aggregateID string) *UserProfileUpdated {
	e := &UserProfileUpdated{}
	e.ID = uuid.New().String()
	e.Source = "/userprofile"
	e.SpecVersion = "1.0"
	e.Type = "com.carddemo.userprofile.updated"
	e.DataContentType = "application/json"
	e.Time = time.Now().Format(time.RFC3339)
	e.Subject = aggregateID
	return e
}

// Type returns the CloudEvent type.
func (e *UserProfileUpdated) Type() string {
	return e.CloudEventEnvelope.Type
}

// AggregateID returns the aggregate ID.
func (e *UserProfileUpdated) AggregateID() string {
	return e.Subject
}
