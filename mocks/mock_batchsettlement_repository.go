package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/batchsettlement/repository"
)

type MockBatchSettlementRepository struct {
	mu   sync.RWMutex
	data map[string]*model.BatchSettlement
}

func NewMockBatchSettlementRepository() *MockBatchSettlementRepository {
	return &MockBatchSettlementRepository{data: make(map[string]*model.BatchSettlement)}
}

var _ repository.BatchSettlementRepository = (*MockBatchSettlementRepository)(nil)

func (m *MockBatchSettlementRepository) Get(id string) (*model.BatchSettlement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

func (m *MockBatchSettlementRepository) Save(aggregate *model.BatchSettlement) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockBatchSettlementRepository) List() ([]*model.BatchSettlement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*model.BatchSettlement, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
