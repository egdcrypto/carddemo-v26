package mocks

import (
	"context"

	"github.com/carddemo/project/src/app/transaction/dto"
)

// MockTransactionService is a mock for the application service layer.
type MockTransactionService struct {
	CreateFunc func(ctx context.Context, req dto.CreateTransactionRequest) (string, error)
	GetFunc     func(ctx context.Context, id string) (*dto.TransactionResponse, error)
	ListFunc    func(ctx context.Context, accountID string, params dto.QueryParams) ([]*dto.TransactionResponse, error)
	VoidFunc    func(ctx context.Context, id string, reason string) error
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (string, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, req)
	}
	return "mock-workflow-id", nil
}

func (m *MockTransactionService) GetTransaction(ctx context.Context, id string) (*dto.TransactionResponse, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return &dto.TransactionResponse{ID: id, Status: "submitted"}, nil
}

func (m *MockTransactionService) ListAccountTransactions(ctx context.Context, accountID string, params dto.QueryParams) ([]*dto.TransactionResponse, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, accountID, params)
	}
	return []*dto.TransactionResponse{}, nil
}

func (m *MockTransactionService) VoidTransaction(ctx context.Context, id string, reason string) error {
	if m.VoidFunc != nil {
		return m.VoidFunc(ctx, id, reason)
	}
	return nil
}
