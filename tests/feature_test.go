package tests

import (
	"testing"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/shared"
)

func TestOpenAccountCmd_Feature(t *testing.T) {
	// Scenario: Successfully execute OpenAccountCmd
	t.Run("success: account opened event emitted", func(t *testing.T) {
		// Given a valid Account aggregate
		agg := model.NewAccount("acc-123")

		// And valid command data
		cmd := command.OpenAccountCmd{
			UserProfileID: "user-101",
			InitialStatus: "Active",
			AccountType:   "Checking",
		}

		// When the command is executed
		events, err := agg.Execute(cmd)

		// Then no error is returned
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}

		// And events are not nil
		if events == nil {
			t.Fatal("expected events list, got nil")
		}

		if len(events) == 0 {
			t.Fatal("expected at least one event, got empty slice")
		}

		// And the event is of type AccountOpenedEvent
		ev, ok := events[0].(*event.AccountOpenedEvent)
		if !ok {
			t.Fatalf("expected AccountOpenedEvent, got %T", events[0])
		}

		// And the event data matches the command
		if ev.UserProfileID != cmd.UserProfileID {
			t.Errorf("expected UserProfileID %s, got %s", cmd.UserProfileID, ev.UserProfileID)
		}
		if ev.Status != cmd.InitialStatus {
			t.Errorf("expected Status %s, got %s", cmd.InitialStatus, ev.Status)
		}
		if ev.AccountType != cmd.AccountType {
			t.Errorf("expected AccountType %s, got %s", cmd.AccountType, ev.AccountType)
		}
	})

	// Scenario: OpenAccountCmd rejected - Account status invariant
	t.Run("rejected: invalid account status", func(t *testing.T) {
		// Given an aggregate command attempting to set an invalid status
		// e.g. trying to Open with 'Suspended' or 'Closed'
		cmd := command.OpenAccountCmd{
			UserProfileID: "user-101",
			InitialStatus: "Suspended", // Violates invariant
			AccountType:   "Checking",
		}

		agg := model.NewAccount("acc-456")

		// When the command is executed
		events, err := agg.Execute(cmd)

		// Then a domain error is expected
		if err == nil {
			t.Fatal("expected error for invalid status, got nil")
		}

		// And error matches the specific invariant violation
		if err != shared.ErrInvalidAccountStatus {
			t.Errorf("expected ErrInvalidAccountStatus, got %v", err)
		}

		// And no events are emitted (failed command)
		if len(events) > 0 {
			t.Errorf("expected no events on failure, got %d", len(events))
		}
	})

	// Scenario: OpenAccountCmd rejected - Account closure is irreversible
	t.Run("rejected: account closure irreversible", func(t *testing.T) {
		// Given an aggregate that attempts to open an account in a closed/invalid state
		// or attempts to modify a closed account
		cmd := command.OpenAccountCmd{
			UserProfileID: "user-101",
			InitialStatus: "Active",
			AccountType:   "Checking",
			IsClosed:      true, // Simulating state that violates the invariant
			CurrentBalance: 100.0,
		}

		agg := model.NewAccount("acc-789")
		// In a real flow, state would be loaded from DB. Here we simulate the check.
		// This tests the business rule enforcement logic.

		// When the command is executed
		// Note: The command payload implies the request was for this state,
		// but the aggregate execution detects the invariant violation.
		events, err := agg.Execute(cmd)

		// Then a domain error is expected
		if err == nil {
			t.Fatal("expected error for closure invariant violation, got nil")
		}

		// And error matches the specific invariant violation
		if err != shared.ErrAccountClosed {
			t.Errorf("expected ErrAccountClosed, got %v", err)
		}

		// And no events are emitted
		if len(events) > 0 {
			t.Errorf("expected no events on failure, got %d", len(events))
		}
	})
}
