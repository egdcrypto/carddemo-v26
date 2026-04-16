package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/event"
)

var (
	ErrInvalidBatchAmount = errors.New("batch amount must be positive")
)

// BatchSettlement represents the Batch Settlement Aggregate.
type BatchSettlement struct {
	shared.AggregateBase
	MerchantID string
	Amount     float64
	Status     string // open, reconciling, completed
}

// NewBatchSettlement creates a new BatchSettlement aggregate.
func NewBatchSettlement(id string) *BatchSettlement {
	return &BatchSettlement{
		AggregateBase: shared.AggregateBase{ID: id},
		Status:        "open",
	}
}

// Execute handles commands.
func (b *BatchSettlement) Execute(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.OpenBatchCmd:
		return b.handleOpenBatch(c)
	case command.ReconcileBatchCmd:
		return b.handleReconcileBatch(c)
	default:
		return errors.New("unknown command")
	}
}

func (b *BatchSettlement) handleOpenBatch(c command.OpenBatchCmd) error {
	if c.Amount <= 0 {
		return ErrInvalidBatchAmount
	}

	b.MerchantID = c.MerchantID
	b.Amount = c.Amount
	b.Status = "reconciling" // Simulate progression

	e := event.BatchOpened{
		BatchID:    b.ID,
		MerchantID: b.MerchantID,
		Amount:     b.Amount,
	}
	b.AddEvent(e)

	return nil
}

func (b *BatchSettlement) handleReconcileBatch(c command.ReconcileBatchCmd) error {
	if b.Status != "reconciling" {
		return errors.New("batch is not in reconciling state")
	}

	b.Status = "completed"
	// No specific event structure provided for Reconcile in existing file, assuming generic update
	return nil
}
