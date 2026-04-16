package model

import "github.com/carddemo/project/src/domain/shared"

const (
	// EventTypeBatchOpened is the event type for a successful batch opening.
	EventTypeBatchOpened = "com.carddemo.batchsettlement.opened"
)

// BatchOpenedEvent is emitted when a new daily batch cycle starts.
type BatchOpenedEvent struct {
	// Meta is handled by the wrapper logic in tests or publishing layer
	BatchID          string `json:"batch_id"`
	SettlementDate   string `json:"settlement_date"`
	OperationalRegion string `json:"operational_region"`
}

// ToDomainEvent converts the struct to the generic DomainEvent interface.
func (e BatchOpenedEvent) ToDomainEvent() shared.DomainEvent {
	return shared.DomainEvent{
		// Type is populated by the publisher or test helper usually, but can be explicit here
		Payload: e,
	}
}
