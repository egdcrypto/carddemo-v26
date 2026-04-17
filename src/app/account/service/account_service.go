package service

import (
	"errors"
	"github.com/google/uuid"

	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/model"
	account_repository "github.com/carddemo/project/src/domain/account/repository"
	userprofile_repository "github.com/carddemo/project/src/domain/userprofile/repository"
)

// AccountService handles the use cases for Account.
type AccountService struct {
	accountRepo account_repository.AccountRepository
	profileRepo userprofile_repository.UserProfileRepository
}

// NewAccountService creates a new service.
func NewAccountService(accountRepo account_repository.AccountRepository, profileRepo userprofile_repository.UserProfileRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		profileRepo: profileRepo,
	}
}

// CreateAccount handles the creation of a new account.
func (s *AccountService) CreateAccount(req dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	// 1. Validate Profile exists (Business Rule consistency)
	_, err := s.profileRepo.Get(req.UserProfileID)
	if err != nil {
		return nil, errors.New("user profile not found")
	}

	// 2. Create Aggregate
	newID := uuid.New().String()
	agg := model.NewAccount(newID)

	// 3. Prepare Command
	cmd := command.OpenAccountCmd{
		UserProfileID: req.UserProfileID,
		InitialStatus: req.Status,
		AccountType:   req.AccountType,
	}

	// 4. Execute
	err = agg.Handle(cmd)
	if err != nil {
		return nil, err
	}

	// 5. Persist
	err = s.accountRepo.Save(agg)
	if err != nil {
		return nil, err
	}

	// 6. Map Response
	return s.mapToResponse(agg), nil
}

// GetAccount retrieves an account by ID.
func (s *AccountService) GetAccount(id string) (*dto.AccountResponse, error) {
	agg, err := s.accountRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if agg == nil {
		return nil, errors.New("account not found")
	}

	return s.mapToResponse(agg), nil
}

// UpdateStatus updates the account status.
func (s *AccountService) UpdateStatus(id string, req dto.UpdateAccountStatusRequest) (*dto.AccountResponse, error) {
	// 1. Load Aggregate
	agg, err := s.accountRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if agg == nil {
		return nil, errors.New("account not found")
	}

	// 2. Execute Command
	cmd := command.UpdateAccountStatusCmd{
		NewStatus: req.NewStatus,
		Reason:    req.Reason,
	}

	err = agg.Handle(cmd)
	if err != nil {
		return nil, err
	}

	// 3. Persist
	err = s.accountRepo.Save(agg)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(agg), nil
}

// DeleteAccount deletes an account.
func (s *AccountService) DeleteAccount(id string) error {
	// In a pure domain model, we might load and mark for deletion.
	// For this REST API, we rely on the repo Delete capability.
	return s.accountRepo.Delete(id)
}

func (s *AccountService) mapToResponse(agg *model.Account) *dto.AccountResponse {
	return &dto.AccountResponse{
		ID:            agg.ID,
		UserProfileID: agg.UserProfileID,
		AccountType:   agg.AccountType,
		Status:        agg.Status,
		Version:       agg.Version,
	}
}
