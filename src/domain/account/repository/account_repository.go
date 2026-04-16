package repository

import (
	"github.com/carddemo/project/src/domain/account/model"
)

// AccountRepository defines the port for account persistence.
type AccountRepository interface {
	Get(id string) (*model.Account, error)
	Save(aggregate *model.Account) error
	Delete(id string) error
	List() ([]*model.Account, error)
}
