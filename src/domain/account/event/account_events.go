package event

import (
	"github.com/carddemo/project/src/domain/shared"
)

// AccountOpenedEvent is emitted when a new account is successfully created.
type AccountOpenedEvent struct {
	shared.DomainEvent
	AccountID     string
	UserProfileID string
	Status        string
	AccountType   string
}
