package rest

import (
	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/cardpolicy/model"
	cardrepo "github.com/carddemo/project/src/domain/card/repository"
	policyrepo "github.com/carddemo/project/src/domain/cardpolicy/repository"
	accountrepo "github.com/carddemo/project/src/domain/account/repository"
	"context"
)

// CardRepository defines the storage interface for Cards used by the port
type CardRepository interface {
	Get(ctx context.Context, id string) (*model.Card, error)
	Save(ctx context.Context, aggregate *model.Card) error
}

// CardPolicyRepository defines the storage interface for Policies used by the port
type CardPolicyRepository interface {
	Get(ctx context.Context, id string) (*model.CardPolicy, error)
	List(ctx context.Context) ([]*model.CardPolicy, error)
	Save(ctx context.Context, aggregate *model.CardPolicy) error
}

// AccountRepository defines the minimal interface needed to validate ownership
type AccountRepository interface {
	Get(ctx context.Context, id string) (interface{}, error) // We just need existence check
}
