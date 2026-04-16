package dto

import (
	"github.com/carddemo/project/src/domain/userprofile/model"
)

// MapToUserProfileResponse converts a domain UserProfile model to a DTO.
func MapToUserProfileResponse(agg *model.UserProfile) UserProfileResponse {
	return UserProfileResponse{
		ID:        agg.ID,
		FirstName: agg.FirstName,
		LastName:  agg.LastName,
		Email:     agg.Email,
		AccountID: agg.AccountID,
	}
}
