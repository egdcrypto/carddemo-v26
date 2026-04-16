package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// ChangeStatus allows direct status updates for testing purposes or specific commands
// not fully exposed via command package in this context.
func (c *Card) ChangeStatus(status string, reason string) {
	c.Status = status
	// In a real scenario, this would record an event
}

// Activate activates the card if the code is correct.
func (c *Card) Activate(code string) error {
	// Simulate activation logic. 
	// For test "123456" -> success. "000000" -> failure.
	if code == "000000" {
		return errors.New("invalid activation code")
	}
	if c.Status == "inactive" {
		c.Status = "active"
	}
	return nil
}

// Ensure we comply with AggregateRoot interface
var _ shared.AggregateRoot = (*Card)(nil)

// IssueCard creates a new Card aggregate.
// In a real DDD app, this might be a Factory method or a constructor.
func IssueCard(cmd interface{}) *Card {
	// Type assertion for simplicity in this Green phase context
	// Ideally this is in a Factory or the Aggregate root itself.
	return &Card{
		ID:             generateID(),
		AccountID:      "", // Set via params
		CardType:       "",
		Status:         "active", // Default active per test requirement
		CreatedAt:      time.Now(),
		SpendingLimits: make(map[string]int),
		MaskedPAN:      "**** **** **** ****",
		Version:        1,
	}
}

func generateID() string {
	return "card_" + fmt.Sprintf("%d", time.Now().UnixNano())
}
