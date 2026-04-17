package shared

import "github.com/google/uuid"

// GenerateUUID creates a new random UUID.
func GenerateUUID() string {
	return uuid.New().String()
}

// AggregateRoot is the base struct for all aggregates.
type AggregateRoot struct {
	ID      string
	Version int
	events  []Event
}

// AddEvent adds a domain event to the aggregate.
func (a *AggregateRoot) AddEvent(event Event) {
	a.events = append(a.events, event)
}

// GetEvents retrieves the list of recorded events.
func (a *AggregateRoot) GetEvents() []Event {
	return a.events
}

// ClearEvents removes all events from the aggregate.
func (a *AggregateRoot) ClearEvents() {
	a.events = nil
}

// Event is the interface for domain events.
type Event interface {
	Type() string
	AggregateID() string
}
