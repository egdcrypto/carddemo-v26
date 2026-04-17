package rest

import (
	"github.com/carddemo/project/src/app/account/handler"
	"github.com/carddemo/project/src/app/account/service"
	userprofile_service "github.com/carddemo/project/src/app/userprofile/service"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/userprofile/repository"
)

// NewAccountHandler wires up dependencies and returns the HTTP handler.
// This acts as the Adapter/Router factory.
func NewAccountHandler(accRepo repository.AccountRepository, profileRepo repository.UserProfileRepository) *handler.AccountHandler {
	accSvc := service.NewAccountService(accRepo, profileRepo)
	profileSvc := userprofile_service.NewUserProfileService(profileRepo)

	return handler.NewAccountHandler(accSvc, profileSvc)
}
