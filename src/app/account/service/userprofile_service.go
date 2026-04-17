package service

import (
	"errors"

	"github.com/carddemo/project/src/domain/userprofile/command"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
	"github.com/google/uuid"
)

// UserProfileApplicationService handles profile use cases.
type UserProfileApplicationService struct {
	repo repository.UserProfileRepository
}

// NewUserProfileApplicationService creates a new service.
func NewUserProfileApplicationService(repo repository.UserProfileRepository) (*UserProfileApplicationService, error) {
	if repo == nil {
		return nil, errors.New("repository required")
	}
	return &UserProfileApplicationService{repo: repo}, nil
}

// LinkOrUpdateProfile links or updates a user profile for an account.
func (s *UserProfileApplicationService) LinkOrUpdateProfile(accountID, firstName, lastName string) (*model.UserProfile, error) {
	// Check if profile exists
	existing, _ := s.repo.GetByAccountID(accountID)

	var agg *model.UserProfile
	if existing != nil {
		agg = existing
	} else {
		id := uuid.New().String()
		agg = model.NewUserProfile(id, accountID, firstName, lastName, "")
	}

	cmd := &command.LinkUserToAccountCommand{
		FirstName: firstName,
		LastName:  lastName,
		AccountID: accountID,
	}

	if err := agg.Execute(cmd); err != nil {
		return nil, err
	}

	if err := s.repo.Save(agg); err != nil {
		return nil, err
	}

	return agg, nil
}
