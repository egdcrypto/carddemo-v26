package model

import (
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	// ErrBatchAlreadyFinalized indicates the batch is in a terminal state.
	ErrBatchAlreadyFinalized = errors.New("batch cannot be modified: already finalized")

	// ErrBalanceMismatch indicates debits and credits do not match.
	ErrBalanceMismatch = errors.New("total debits must equal total credits to zero out settlement balance")

	// ErrPendingTransactionsExist indicates the batch cannot be finalized.
	ErrPendingTransactionsExist = errors.New("settlement batch cannot be finalized: contains uncommitted or pending transactions")
)

// State represents the lifecycle state of the batch.
type State string

const (
	StateOpen       State = "OPEN"
	StateFinalized  State = "FINALIZED"
	StateRejected   State = "REJECTED"
)

// BatchSettlement represents the BatchSettlement aggregate.
type BatchSettlement struct {
	shared.AggregateRoot
	ID                  string
	State               State
	SettlementDate      string
	OperationalRegion   string
	TotalDebits         int64 // in cents
	TotalCredits        int64 // in cents
	HasPendingTxns      bool
	OpenedAt            *time.Time
}

// NewBatchSettlement creates a new BatchSettlement instance.
func NewBatchSettlement(id string) *BatchSettlement {
	return &BatchSettlement{
		ID:    id,
		State: StateOpen, // Default assumption for a new object before explicit command, or empty
	}
}

// BatchSettlementInState creates a batch for testing purposes with a specific state setup.
// This helper facilitates testing invariants without complex event history hydration.
func BatchSettlementInState(id string, opts ...BatchOption) *BatchSettlement {
	b := &BatchSettlement{
		ID:    id,
		State: StateOpen,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// BatchOption configures a BatchSettlement.
type BatchOption func(*BatchSettlement)

// StateWithPendingTxns sets the batch to have pending transactions.
func StateWithPendingTxns(b *BatchSettlement) {
	b.HasPendingTxns = true
}

// StateWithBalanceMismatch sets the batch to have unequal debits/credits.
func StateWithBalanceMismatch(b *BatchSettlement) {
	b.TotalDebits = 100
	b.TotalCredits = 0
}

// Execute handles commands for the BatchSettlement aggregate.
func (b *BatchSettlement) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.OpenBatchCmd:
		return b.handleOpenBatch(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleOpenBatch processes the OpenBatchCmd.
// TDD: Initially returns nil, nil to fail the "event emitted" check, then evolves.
func (b *BatchSettlement) handleOpenBatch(cmd command.OpenBatchCmd) ([]shared.DomainEvent, error) {
	return nil, nil
}

// ID returns the aggregate ID.
func (b *BatchSettlement) GetID() string {
	return b.ID
}
