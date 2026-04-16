package adapter

import (
	"github.com/carddemo/project/src/app/transaction/dto"
	"github.com/carddemo/project/src/domain/transaction/model"
)

// ToTransactionResponse maps a Transaction model to a DTO.
func ToTransactionResponse(agg *model.Transaction) *dto.TransactionResponse {
	if agg == nil {
		return nil
	}
	return &dto.TransactionResponse{
		ID:              agg.ID,
		AccountID:       agg.AccountID,
		CardID:          agg.CardID,
		Amount:          agg.Amount,
		TransactionType: agg.TransactionType,
		Status:          agg.Status,
		CreatedAt:       agg.CreatedAt,
	}
}
