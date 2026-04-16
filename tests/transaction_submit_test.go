package tests

import (
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/event"
	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/stretchr/testify/assert"
)

// TDD Red Phase Tests for SubmitTransactionCmd
// These tests expect the domain logic to be implemented in model/transaction_aggregate.go

func TestSubmitTransactionCmd_Success(t *testing.T) {
	// Given
	tx := model.NewTransaction("tx-123")
	cmd := command.SubmitTransactionCmd{
		TransactionID: "tx-123",
		AccountID:     "acc-456",
		CardID:        "card-789",
		Amount:        100.50,
		Type:          "purchase",
		AccountStatus: "Active",
	}

	// When
	events, err := tx.Execute(cmd)

	// Then
	assert.NoError(t, err, "Execute should not return an error for valid data")
	assert.NotNil(t, events, "Events should not be nil")
	assert.Len(t, events, 1, "Should emit exactly one event")

	// Verify specific Event Type and Payload
	ev, ok := events[0].(*event.TransactionSubmitted)
	assert.True(t, ok, "Event should be of type TransactionSubmitted")
	assert.Equal(t, "tx-123", ev.AggregateID)
	assert.Equal(t, "acc-456", ev.AccountID)
	assert.Equal(t, 100.50, ev.Amount)
	assert.Equal(t, "purchase", ev.Type)

	// Allow slight tolerance for timestamp check
	assert.WithinDuration(t, time.Now(), ev.OccurredAt, 2*time.Second)
}

func TestSubmitTransactionCmd_Rejected_InvalidAmount(t *testing.T) {
	// Scenario: Transaction amount must be strictly greater than zero
	// Given
	tx := model.NewTransaction("tx-invalid-amt")
	cmd := command.SubmitTransactionCmd{
		TransactionID: "tx-invalid-amt",
		AccountID:     "acc-valid",
		CardID:        "card-valid",
		Amount:        0, // Violation: Amount <= 0
		Type:          "purchase",
		AccountStatus: "Active",
	}

	// When
	events, err := tx.Execute(cmd)

	// Then
	assert.Error(t, err, "Execute should return an error for non-positive amount")
	assert.Nil(t, events, "No events should be emitted on command rejection")

	// Checking for specific domain error validation (Implementation detail optional for Red phase, but good practice)
	// Assuming a typed error exists in shared package or handling string match
	// assert.ErrorIs(t, err, shared.ErrInvalidAmount)
}

func TestSubmitTransactionCmd_Rejected_AccountNotActive(t *testing.T) {
	// Scenario: Account must be in 'Active' status
	// Given
	tx := model.NewTransaction("tx-inactive-acc")
	cmd := command.SubmitTransactionCmd{
		TransactionID: "tx-inactive-acc",
		AccountID:     "acc-suspended",
		CardID:        "card-valid",
		Amount:        50.00,
		Type:          "purchase",
		AccountStatus: "Suspended", // Violation: Not 'Active'
	}

	// When
	events, err := tx.Execute(cmd)

	// Then
	assert.Error(t, err, "Execute should return an error when account status is not Active")
	assert.Nil(t, events, "No events should be emitted on command rejection")
	assert.Contains(t, err.Error(), "Active", "Error message should reference the Active status requirement")
}
