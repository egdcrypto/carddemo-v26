package service

import (
	"errors"
	"fmt"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
)

var (
	// ErrAccountNotFound is returned when an account cannot be found in the repository.
	ErrAccountNotFound = errors.New("account not found")
)

// AccountService is the application service for account use cases.
// It acts as the orchestrator, loading aggregates, dispatching commands, and persisting state.
type AccountService struct {
	repo repository.AccountRepository
}

// NewAccountService creates a new AccountService.
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

// CreateAccount handles the creation of a new account.
func (s *AccountService) CreateAccount(cmd command.OpenAccountCmd) (*model.Account, error) {
	// In a real scenario, we might verify the UserProfileID exists via a UserProfileService.
	// For now, we assume it exists (or the aggregate handles validation).

	// Create the aggregate using the factory method.
	// Ideally, IDs are generated here or in the factory.
	// Assuming the Aggregate factory handles ID generation if empty.
	agg := model.NewAccount(cmd)

	if err := s.repo.Save(agg); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	return agg, nil
}

// GetAccount retrieves an account by ID.
func (s *AccountService) GetAccount(id string) (*model.Account, error) {
	agg, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if agg == nil {
		return nil, ErrAccountNotFound
	}
	return agg, nil
}

// UpdateAccountStatus handles status updates for an account.
func (s *AccountService) UpdateAccountStatus(id string, cmd command.UpdateAccountStatusCmd) (*model.Account, error) {
	// 1. Load the aggregate
	agg, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load account: %w", err)
	}
	if agg == nil {
		return nil, ErrAccountNotFound
	}

	// 2. Execute the command on the aggregate
	// This updates the internal state and generates domain events
	if err := agg.Execute(cmd); err != nil {
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	// 3. Persist the changes
	if err := s.repo.Save(agg); err != nil {
		return nil, fmt.Errorf("failed to save updated account: %w", err)
	}

	return agg, nil
}

// DeleteAccount handles the deletion of an account.
func (s *AccountService) DeleteAccount(id string) error {
	// Check if exists first to return correct 404 logic if needed, 
	// though Delete operation usually is idempotent.
	_, err := s.repo.Get(id)
	if err != nil {
		return fmt.Errorf("failed to check account existence: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}
