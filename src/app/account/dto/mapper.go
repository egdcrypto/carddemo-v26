package dto

import (
	"github.com/carddemo/project/src/domain/account/model"
)

// MapToAccountResponse converts a domain Account model to a DTO.
func MapToAccountResponse(agg *model.Account) AccountResponse {
	return AccountResponse{
		ID:            agg.ID,
		UserProfileID: agg.UserProfileID,
		Status:        agg.Status,
		AccountType:   agg.AccountType,
		CreatedAt:     agg.CreatedAt,
		UpdatedAt:     agg.UpdatedAt,
	}
}
