package model

import (
	"github.com/carddemo/project/src/domain/cardpolicy/command"
)

// UpdateLimits handles the command to update limits.
func (cp *CardPolicy) UpdateLimits(cmd command.UpdateCardLimitsCmd) {
	cp.DailyLimit = cmd.DailyLimit
	cp.WeeklyLimit = cmd.WeeklyLimit
	// In real app: check invariants, record event
}

// AssignCardPolicy links a policy to a card (implied by repository usage).
func AssignCardPolicy(cardID string) *CardPolicy {
	return &CardPolicy{
		CardID:     cardID,
		DailyLimit: 0,
		WeeklyLimit: 0,
		IsActive:   true,
	}
}
