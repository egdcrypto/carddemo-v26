package mocks

import (
	"context"

	"github.com/carddemo/project/src/app/batchsettlement/dto"
)

// MockBatchSettlementService is a mock for the batch settlement application service.
type MockBatchSettlementService struct {
	CreateFunc func(ctx context.Context, req dto.CreateBatchSettlementRequest) (string, error)
	GetFunc     func(ctx context.Context, id string) (*dto.BatchSettlementResponse, error)
	ListFunc    func(ctx context.Context) ([]*dto.BatchSettlementResponse, error)
}

func (m *MockBatchSettlementService) CreateSettlement(ctx context.Context, req dto.CreateBatchSettlementRequest) (string, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, req)
	}
	return "settlement-123", nil
}

func (m *MockBatchSettlementService) GetSettlement(ctx context.Context, id string) (*dto.BatchSettlementResponse, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return &dto.BatchSettlementResponse{ID: id, Status: "open"}, nil
}

func (m *MockBatchSettlementService) ListSettlements(ctx context.Context) ([]*dto.BatchSettlementResponse, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []*dto.BatchSettlementResponse{}, nil
}
