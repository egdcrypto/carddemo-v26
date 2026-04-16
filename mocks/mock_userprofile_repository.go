package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
)

// MockUserProfileRepository is an in-memory implementation of repository.UserProfileRepository.
type MockUserProfileRepository struct {
	mu   sync.RWMutex
	data map[string]*model.UserProfile
}

// NewMockUserProfileRepository creates a new mock repository.
func NewMockUserProfileRepository() *MockUserProfileRepository {
	return &MockUserProfileRepository{
		data: make(map[string]*model.UserProfile),
	}
}

// Ensure MockUserProfileRepository implements the interface.
var _ repository.UserProfileRepository = (*MockUserProfileRepository)(nil)

// Get retrieves an aggregate by ID.
func (m *MockUserProfileRepository) Get(id string) (*model.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

// GetByAccountID retrieves a profile by the associated account ID.
func (m *MockUserProfileRepository) GetByAccountID(accountID string) (*model.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.data {
		// Simplified check, assumes model.UserProfile has an AccountID field or relationship
		// For test red phase, we return nil if not found explicitly set up
		if v.AccountID == accountID {
			return v, nil
		}
	}
	return nil, nil // or return error not found
}

// Save stores an aggregate.
func (m *MockUserProfileRepository) Save(aggregate *model.UserProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

// Delete removes an aggregate.
func (m *MockUserProfileRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}
