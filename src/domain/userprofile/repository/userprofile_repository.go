package repository

import (
	userprofile_model "github.com/carddemo/project/src/domain/userprofile/model"
)

// UserProfileRepository defines the persistence interface for User Profiles.
type UserProfileRepository interface {
	Get(id string) (*userprofile_model.UserProfile, error)
	GetByAccountID(accountID string) (*userprofile_model.UserProfile, error)
	Save(aggregate *userprofile_model.UserProfile) error
	Delete(id string) error
	List() ([]*userprofile_model.UserProfile, error)
}
