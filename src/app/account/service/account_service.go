package service

import (
	"errors"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/google/uuid"
)

var (
	ErrRepoRequired = errors.New("repository required")
)

// AccountApplicationService handles use cases.
type AccountApplicationService struct {
	repo repository.AccountRepository
}

// NewAccountApplicationService creates a new service.
func NewAccountApplicationService(repo repository.AccountRepository) (*AccountApplicationService, error) {
	if repo == nil {
		return nil, ErrRepoRequired
	}
	return &AccountApplicationService{repo: repo}, nil
}

// CreateAccount creates a new account.
func (s *AccountApplicationService) CreateAccount userProfileID, accountType, status string) (*model.Account, error) {
	id := uuid.New().String()
	agg := model.NewAccount(id, userProfileID, status, accountType)

	cmd := &command.OpenAccountCmd{
		UserProfileID: userProfileID,
		InitialStatus: status,
		AccountType:   accountType,
	}

	if err := agg.Execute(cmd); err != nil {
		return nil, err
	}

	if err := s.repo.Save(agg); err != nil {
		return nil, err
	}

	return agg, nil
}

// GetAccount retrieves an account by ID.
func (s *AccountApplicationService) GetAccount(id string) (*model.Account, error) {
	return s.repo.Get(id)
}

// UpdateAccountStatus updates the status of an account.
func (s *AccountApplicationService) UpdateAccountStatus(id, newStatus, reason string) error {
	agg, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	if agg == nil {
		return model.ErrAccountNotFound
	}

	cmd := &command.UpdateAccountStatusCmd{
		NewStatus: newStatus,
		Reason:    reason,
	}

	if err := agg.Execute(cmd); err != nil {
		return err
	}

	return s.repo.Save(agg)
}

// DeleteAccount deletes an account.
func (s *AccountApplicationService) DeleteAccount(id string) error {
	return s.repo.Delete(id)
}
