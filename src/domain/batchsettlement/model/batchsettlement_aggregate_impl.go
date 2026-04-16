package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/batchsettlement/command"
)

// Handle executes commands on the BatchSettlement aggregate.
func (b *BatchSettlement) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.OpenBatchCmd:
		return b.handleOpen(c)
	default:
		return errors.New("unknown command")
	}
}

func (b *BatchSettlement) handleOpen(cmd command.OpenBatchCmd) error {
	// Domain Validation
	if len(cmd.Currency) != 3 {
		return errors.New("invalid currency format")
	}

	// Apply state changes
	b.Status = "Open"
	return nil
}
