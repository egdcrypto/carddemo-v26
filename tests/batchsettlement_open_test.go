package tests

import (
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOpenBatchCmd_Success tests the successful execution of the command.
func TestOpenBatchCmd_Success(t *testing.T) {
	// Arrange
	batchID := "batch-123"
	batch := model.NewBatchSettlement(batchID)

	cmd := command.OpenBatchCmd{
		SettlementDate:    "2023-10-27",
		OperationalRegion: "US-EAST",
	}

	// Act
	events, err := batch.Execute(cmd)

	// Assert
	require.NoError(t, err, "Command should execute successfully")
	require.NotNil(t, events, "Events should not be nil")
	assert.Len(t, events, 1, "Exactly one event should be emitted")

	// Verify Event Content (TDD Check)
	// We expect the payload to be BatchOpenedEvent.
	// Since the implementation is currently empty/null return, this checks the structure.
	castedEvent, ok := events[0].Payload.(model.BatchOpenedEvent)
	require.True(t, ok, "Event payload should be BatchOpenedEvent")

	assert.Equal(t, batchID, castedEvent.BatchID)
	assert.Equal(t, "2023-10-27", castedEvent.SettlementDate)
	assert.Equal(t, "US-EAST", castedEvent.OperationalRegion)

	// Verify State Changes (Implicit in Red Phase, but good to define future expectation)
	// assert.Equal(t, model.StateOpen, batch.State) // This will pass/fail depending on impl
}

// TestOpenBatchCmd_Rejected_PendingTransactions tests rejection due to pending transactions.
func TestOpenBatchCmd_Rejected_PendingTransactions(t *testing.T) {
	// Arrange
	batchID := "batch-456"
	// Use the helper to set state that violates the invariant
	// In this scenario, we are simulating an attempt to open/finalize a batch
	// that logic dictates cannot proceed because of existing pending items.
	batch := model.BatchSettlementInState(batchID, model.StateWithPendingTxns)

	cmd := command.OpenBatchCmd{
		SettlementDate:    "2023-10-27",
		OperationalRegion: "US-WEST",
	}

	// Act
	events, err := batch.Execute(cmd)

	// Assert
	require.Error(t, err, "Command should be rejected")
	assert.Nil(t, events, "No events should be emitted on rejection")
	assert.Equal(t, model.ErrPendingTransactionsExist, err)
}

// TestOpenBatchCmd_Rejected_BalanceMismatch tests rejection due to imbalance.
func TestOpenBatchCmd_Rejected_BalanceMismatch(t *testing.T) {
	// Arrange
	batchID := "batch-789"
	// Use the helper to set state that violates the invariant
	// (Debits != Credits)
	batch := model.BatchSettlementInState(batchID, model.StateWithBalanceMismatch)

	cmd := command.OpenBatchCmd{
		SettlementDate:    "2023-10-27",
		OperationalRegion: "EU-CENTRAL",
	}

	// Act
	events, err := batch.Execute(cmd)

	// Assert
	require.Error(t, err, "Command should be rejected")
	assert.Nil(t, events, "No events should be emitted on rejection")
	assert.Equal(t, model.ErrBalanceMismatch, err)
}
