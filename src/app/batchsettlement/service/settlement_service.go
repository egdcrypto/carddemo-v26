package service

import (
	"fmt"
	"time"

	"github.com/carddemo/project/src/app/batchsettlement/adapter"
	"github.com/carddemo/project/src/app/batchsettlement/dto"
	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/batchsettlement/repository"
)

// SettlementApplication handles the use cases for batch settlements.
type SettlementApplication struct {
	repo repository.BatchSettlementRepository
}

// NewSettlementApplication creates a new SettlementApplication.
func NewSettlementApplication(repo repository.BatchSettlementRepository) *SettlementApplication {
	return &SettlementApplication{repo: repo}
}

// Create creates a new batch settlement.
func (s *SettlementApplication) Create(req dto.CreateSettlementRequest) (*dto.SettlementResponse, error) {
	// Generate ID
	id := fmt.Sprintf("batch_%s_%d", req.MerchantID, time.Now().Unix())

	// 1. Load Aggregate
	aggregate := model.NewBatchSettlement(id, req.MerchantID, req.Currency)

	// 2. Prepare Command
	cmd := command.OpenBatchCmd{
		BatchID:    id,
		MerchantID: req.MerchantID,
		Currency:   req.Currency,
	}

	// 3. Execute on Aggregate
	err := aggregate.Handle(cmd)
	if err != nil {
		return nil, err
	}

	// 4. Persist
	if err := s.repo.Save(aggregate); err != nil {
		return nil, err
	}

	// 5. Map to Response
	return adapter.ToSettlementResponse(aggregate), nil
}

// Get retrieves a settlement by ID.
func (s *SettlementApplication) Get(id string) (*dto.SettlementResponse, error) {
	aggregate, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return adapter.ToSettlementResponse(aggregate), nil
}

// List retrieves all settlements.
func (s *SettlementApplication) List() ([]*dto.SettlementResponse, error) {
	aggregates, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	var result []*dto.SettlementResponse
	for _, agg := range aggregates {
		result = append(result, adapter.ToSettlementResponse(agg))
	}
	return result, nil
}
