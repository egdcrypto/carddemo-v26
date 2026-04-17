package model

import (
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrInvalidAccountType   = errors.New("invalid account type")
	ErrOptimisticLockFailed = errors.New("version mismatch")
)

// Account is the aggregate root.
type Account struct {
	shared.AggregateRoot
	ID            string
	UserProfileID  string
	Status        string
	AccountType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewAccount creates a new Account aggregate.
func NewAccount(id, userProfileID, status, accountType string) *Account {
	return &Account{
		AggregateRoot: shared.AggregateRoot{Version: 0},
		ID:            id,
		UserProfileID: userProfileID,
		Status:        status,
		AccountType:   accountType,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// Execute handles commands.
func (a *Account) Execute(cmd interface{}) error {
	switch c := cmd.(type) {
	case *command.OpenAccountCmd:
		return a.openAccount(c)
	case *command.UpdateAccountStatusCmd:
		return a.updateStatus(c)
	default:
		return shared.ErrUnknownCommand
	}
}

func (a *Account) openAccount(cmd *command.OpenAccountCmd) error {
	// In a real scenario, we might validate transitions or business rules here.
	// For this phase, we assume the aggregate is new or state management is handled externally.
	a.UserProfileID = cmd.UserProfileID
	a.Status = cmd.InitialStatus
	a.AccountType = cmd.AccountType
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	evt := event.NewAccountOpened(a.ID, cmd)
	evt.Payload.AccountID = a.ID
	evt.Payload.UserProfileID = a.UserProfileID
	evt.Payload.Status = a.Status
	evt.Payload.AccountType = a.AccountType

	a.AddEvent(evt)
	return nil
}

func (a *Account) updateStatus(cmd *command.UpdateAccountStatusCmd) error {
	oldStatus := a.Status
	a.Status = cmd.NewStatus
	a.UpdatedAt = time.Now()

	evt := event.NewAccountStatusUpdated(a.ID)
	evt.Payload.AccountID = a.ID
	evt.Payload.OldStatus = oldStatus
	evt.Payload.NewStatus = cmd.NewStatus
	evt.Payload.Reason = cmd.Reason

	a.AddEvent(evt)
	return nil
}
