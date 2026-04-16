package mocks

import (
	"context"
	"errors"
	"sync"
)

// MockAccountRepository is a mock for AccountRepository
type MockAccountRepository struct {
	mu   sync.RWMutex
	data map[string]interface{} // stored as empty struct or mock
}

func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{data: make(map[string]interface{})}
}

func (m *MockAccountRepository) Get(ctx context.Context, id string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.data[id]; !ok {
		return nil, errors.New("account not found")
	}
	return struct{ ID string }{ID: id}, nil
}
