package mocks

import (
	"sync"

	userprofile_model "github.com/carddemo/project/src/domain/userprofile/model"
	userprofile_repository "github.com/carddemo/project/src/domain/userprofile/repository"
)

// MockUserProfileRepository is an in-memory implementation for testing.
type MockUserProfileRepository struct {
	mu   sync.RWMutex
	data map[string]*userprofile_model.UserProfile
}

func NewMockUserProfileRepository() *MockUserProfileRepository {
	return &MockUserProfileRepository{data: make(map[string]*userprofile_model.UserProfile)}
}

var _ userprofile_repository.UserProfileRepository = (*MockUserProfileRepository)(nil)

func (m *MockUserProfileRepository) Get(id string) (*userprofile_model.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if val, ok := m.data[id]; ok {
		return val, nil
	}
	return nil, nil // Not found
}

func (m *MockUserProfileRepository) GetByAccountID(accountID string) (*userprofile_model.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.data {
		if v.AccountID == accountID {
			return v, nil
		}
	}
	return nil, nil
}

func (m *MockUserProfileRepository) Save(aggregate *userprofile_model.UserProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockUserProfileRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MockUserProfileRepository) List() ([]*userprofile_model.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*userprofile_model.UserProfile, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
