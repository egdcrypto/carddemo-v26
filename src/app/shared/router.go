package rest

import (
	"github.com/carddemo/project/src/domain/account/repository"
	userprofilerepository "github.com/carddemo/project/src/domain/userprofile/repository"
	"github.com/go-chi/chi/v5"
)

// WireAccountRoutes connects the handlers to the router using the provided repositories.
// This acts as the composition root for the REST layer.
func WireAccountRoutes(r *chi.Mux, accountRepo repository.AccountRepository, userProfileRepo userprofilerepository.UserProfileRepository) {
	// Initialize Application Services
	accountService := NewAccountService(accountRepo)
	userProfileService := NewUserProfileService(userProfileRepo)

	// Initialize Handlers
	accHandler := NewAccountHandler(accountService)
	profileHandler := NewUserProfileHandler(userProfileService)

	// Register Routes
	// Note: We mount both. AccountHandler handles /accounts root.
	// UserProfileHandler handles /accounts/{id}/profile.
	// To avoid duplicate route definitions in chi, we can chain them or
	// let one handler manage the mounting. For simplicity here, we call both.

	accHandler.RegisterRoutes(r)
	profileHandler.RegisterRoutes(r)
}
