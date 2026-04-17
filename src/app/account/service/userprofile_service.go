package service

import (
	"errors"
	"github.com/google/uuid"

	"github.com/carddemo/project/src/app/userprofile/dto"
	"github.com/carddemo/project/src/domain/userprofile/command"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
)

// UserProfileService handles use cases for UserProfile.
type UserProfileService struct {
	profileRepo repository.UserProfileRepository
}

// NewUserProfileService creates a new service.
func NewUserProfileService(profileRepo repository.UserProfileRepository) *UserProfileService {
	return &UserProfileService{
		profileRepo: profileRepo,
	}
}

// GetProfileByAccountID retrieves a profile by its associated account ID.
func (s *UserProfileService) GetProfileByAccountID(accountID string) (*dto.UserProfileResponse, error) {
	profile, err := s.profileRepo.GetByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found for account")
	}

	return s.mapToResponse(profile), nil
}

// UpdateProfileByAccountID updates a profile identified by account ID.
func (s *UserProfileService) UpdateProfileByAccountID(accountID string, req dto.UpdateUserProfileRequest) (*dto.UserProfileResponse, error) {
	// 1. Load Aggregate
	profile, err := s.profileRepo.GetByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		// If profile doesn't exist for the account, we might create it or return 404.
		// Assuming 404 based on typical resource semantics.
		return nil, errors.New("profile not found for account")
	}

	// 2. Execute Command
	cmd := command.UpdateProfileCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	err = profile.Handle(cmd)
	if err != nil {
		return nil, err
	}

	// 3. Persist
	err = s.profileRepo.Save(profile)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(profile), nil
}

// CreateProfile helper (not used in tests but good for completeness)
func (s *UserProfileService) CreateProfile(accountID string, req dto.UpdateUserProfileRequest) (*dto.UserProfileResponse, error) {
	newID := uuid.New().String()
	profile := model.NewUserProfile(newID, accountID)

	cmd := command.UpdateProfileCommand{ // Reuse update command structure as fields match
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	err := profile.Handle(cmd)
	if err != nil {
		return nil, err
	}

	err = s.profileRepo.Save(profile)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(profile), nil
}

func (s *UserProfileService) mapToResponse(p *model.UserProfile) *dto.UserProfileResponse {
	return &dto.UserProfileResponse{
		ID:        p.ID,
		AccountID: p.AccountID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
	}
}
