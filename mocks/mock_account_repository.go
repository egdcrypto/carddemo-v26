package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
)

// MockAccountRepository is an in-memory implementation of AccountRepository for testing.
type MockAccountRepository struct {
	mu   sync.RWMutex
	data map[string]*model.Account
}

// NewMockAccountRepository creates a new empty mock repository.
func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{data: make(map[string]*model.Account)}
}

var _ repository.AccountRepository = (*MockAccountRepository)(nil)

// Get retrieves an account by ID.
func (m *MockAccountRepository) Get(id string) (*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

// Save saves an account.
func (m *MockAccountRepository) Save(aggregate *model.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

// Delete removes an account.
func (m *MockAccountRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

// List returns all accounts.
func (m *MockAccountRepository) List() ([]*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*model.Account, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
