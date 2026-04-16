package repository

import (
	"github.com/carddemo/project/src/domain/batchsettlement/model"
)

// BatchSettlementRepository defines the storage interface for BatchSettlement aggregates.
type BatchSettlementRepository interface {
	Get(id string) (*model.BatchSettlement, error)
	Save(aggregate *model.BatchSettlement) error
	Delete(id string) error
	List() ([]*model.BatchSettlement, error)
}
