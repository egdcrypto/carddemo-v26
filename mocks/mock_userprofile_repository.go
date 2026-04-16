package mocks

import (
	"sync"

	userprofilemodel "github.com/carddemo/project/src/domain/userprofile/model"
	userprofilerepository "github.com/carddemo/project/src/domain/userprofile/repository"
)

// MockUserProfileRepository is an in-memory implementation of repository.UserProfileRepository.
type MockUserProfileRepository struct {
	mu   sync.RWMutex
	data map[string]*userprofilemodel.UserProfile
}

// NewMockUserProfileRepository creates a new mock repository.
func NewMockUserProfileRepository() *MockUserProfileRepository {
	return &MockUserProfileRepository{
		data: make(map[string]*userprofilemodel.UserProfile),
	}
}

// Ensure MockUserProfileRepository implements the interface.
var _ userprofilerepository.UserProfileRepository = (*MockUserProfileRepository)(nil)

// Get retrieves an aggregate by ID.
func (m *MockUserProfileRepository) Get(id string) (*userprofilemodel.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if val, ok := m.data[id]; ok {
		return val, nil
	}
	return nil, nil // Return nil to simulate not found
}

// GetByAccountID retrieves a profile by its linked Account ID.
func (m *MockUserProfileRepository) GetByAccountID(accountID string) (*userprofilemodel.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.data {
		if v.AccountID == accountID {
			return v, nil
		}
	}
	return nil, nil
}

// Save stores an aggregate.
func (m *MockUserProfileRepository) Save(aggregate *userprofilemodel.UserProfile) error {
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

// List returns all aggregates.
func (m *MockUserProfileRepository) List() ([]*userprofilemodel.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*userprofilemodel.UserProfile, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
