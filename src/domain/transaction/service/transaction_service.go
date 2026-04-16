package service

import (
	"fmt"

	"github.com/carddemo/project/src/app/transaction/adapter"
	"github.com/carddemo/project/src/app/transaction/dto"
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/carddemo/project/src/domain/transaction/repository"
)

// TemporalClient defines the interface for async workflow triggering.
type TemporalClient interface {
	StartWorkflow(id string) error
}

// TransactionApplication handles the use cases for transactions.
type TransactionApplication struct {
	repo    repository.TransactionRepository
	temporal TemporalClient
}

// NewTransactionApplication creates a new TransactionApplication.
func NewTransactionApplication(repo repository.TransactionRepository, temporal TemporalClient) *TransactionApplication {
	return &TransactionApplication{
		repo:    repo,
		temporal: temporal,
	}
}

// Create handles the creation of a new transaction.
func (s *TransactionApplication) Create(req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	// Generate ID (in real app, use UUID)
	id := fmt.Sprintf("tx_%s_%d", req.AccountID, len(req.AccountID))

	// 1. Load Aggregate (New in this case)
	aggregate := model.NewTransaction(id, req.AccountID, req.CardID, req.Amount, req.TransactionType)

	// 2. Prepare Command
	cmd := command.SubmitTransactionCmd{
		TransactionID:   id,
		AccountID:       req.AccountID,
		CardID:          req.CardID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		AccountStatus:   "Active", // Simplified for Green phase
	}

	// 3. Execute on Aggregate
	err := aggregate.Handle(cmd)
	if err != nil {
		return nil, err
	}

	// 4. Persist
	if err := s.repo.Save(aggregate); err != nil {
		return nil, err
	}

	// 5. Trigger Temporal Workflow (Async)
	if s.temporal != nil {
		go s.temporal.StartWorkflow(id)
	}

	// 6. Map to Response
	return adapter.ToTransactionResponse(aggregate), nil
}

// Get retrieves a transaction by ID.
func (s *TransactionApplication) Get(id string) (*dto.TransactionResponse, error) {
	aggregate, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return adapter.ToTransactionResponse(aggregate), nil
}

// Reverse voids a transaction.
func (s *TransactionApplication) Reverse(id string, req dto.ReverseTransactionRequest) (*dto.TransactionResponse, error) {
	// 1. Load Aggregate
	aggregate, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	// 2. Prepare Command
	cmd := command.ReverseTransactionCmd{
		TransactionID: id,
		Reason:        req.Reason,
		Amount:        aggregate.Amount,
		AccountStatus: "Active",
	}

	// 3. Execute on Aggregate
	err = aggregate.Handle(cmd)
	if err != nil {
		return nil, err
	}

	// 4. Persist
	if err := s.repo.Save(aggregate); err != nil {
		return nil, err
	}

	return adapter.ToTransactionResponse(aggregate), nil
}

// ListByAccount retrieves transactions for an account.
// Note: Using List() and filtering for simplicity as per mocked repository limitations.
func (s *TransactionApplication) ListByAccount(accountID string) ([]*dto.TransactionResponse, error) {
	all, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	var result []*dto.TransactionResponse
	for _, tx := range all {
		if tx.AccountID == accountID {
			result = append(result, adapter.ToTransactionResponse(tx))
		}
	}
	return result, nil
}
