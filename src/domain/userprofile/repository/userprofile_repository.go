package repository

import (
	"github.com/carddemo/project/src/domain/userprofile/model"
)

// UserProfileRepository defines the port for interacting with UserProfile aggregates.
type UserProfileRepository interface {
	Get(id string) (*model.UserProfile, error)
	GetByAccountID(accountID string) (*model.UserProfile, error)
	Save(aggregate *model.UserProfile) error
	Delete(id string) error
}
