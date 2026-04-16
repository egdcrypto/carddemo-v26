package adapter

import (
	"github.com/carddemo/project/src/app/batchsettlement/dto"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
)

// ToSettlementResponse maps a BatchSettlement model to a DTO.
func ToSettlementResponse(agg *model.BatchSettlement) *dto.SettlementResponse {
	if agg == nil {
		return nil
	}
	return &dto.SettlementResponse{
		ID:         agg.ID,
		MerchantID: agg.MerchantID,
		Currency:   agg.Currency,
		Status:     agg.Status,
		CreatedAt:  agg.CreatedAt,
	}
}
