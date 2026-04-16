package service

import (
	"fmt"

	"github.com/carddemo/project/src/domain/userprofile/command"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
)

// UserProfileService is the application service for user profile use cases.
type UserProfileService struct {
	repo repository.UserProfileRepository
}

// NewUserProfileService creates a new UserProfileService.
func NewUserProfileService(repo repository.UserProfileRepository) *UserProfileService {
	return &UserProfileService{
		repo: repo,
	}
}

// GetProfileByAccountID retrieves a profile associated with a specific account ID.
func (s *UserProfileService) GetProfileByAccountID(accountID string) (*model.UserProfile, error) {
	profile, err := s.repo.GetByAccountID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	if profile == nil {
		return nil, ErrProfileNotFound
	}
	return profile, nil
}

// UpdateProfileDetails updates the details of a user profile.
func (s *UserProfileService) UpdateProfileDetails(id string, cmd command.UpdateProfileDetailsCmd) (*model.UserProfile, error) {
	// 1. Load the aggregate
	agg, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load profile: %w", err)
	}
	if agg == nil {
		return nil, ErrProfileNotFound
	}

	// 2. Execute command
	if err := agg.Execute(cmd); err != nil {
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	// 3. Persist
	if err := s.repo.Save(agg); err != nil {
		return nil, fmt.Errorf("failed to save updated profile: %w", err)
	}

	return agg, nil
}
