package model

import (
	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/userprofile/command"
)

// UserProfile is the aggregate root for the UserProfile domain.
type UserProfile struct {
	shared.AggregateRoot
	AccountID string
	FirstName string
	LastName  string
	Email     string
}

// Handle processes commands for the UserProfile aggregate.
func (p *UserProfile) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.UpdateProfileDetailsCmd:
		return p.updateDetails(c)
	case command.LinkUserToAccountCmd:
		// Handle linking if necessary
		return nil
	default:
		return nil
	}
}

func (p *UserProfile) updateDetails(cmd command.UpdateProfileDetailsCmd) error {
	p.FirstName = cmd.FirstName
	p.LastName = cmd.LastName
	p.Email = cmd.Email
	p.Version++
	return nil
}
