package rest_test

import (
	"github.com/carddemo/project/src/domain/account/repository"
	userprofilerepository "github.com/carddemo/project/src/domain/userprofile/repository"
	"github.com/carddemo/project/src/mocks"
	"github.com/go-chi/chi/v5"
)

// InitializeTestHandlers is a helper function to wire up the REST handlers for testing.
// It abstracts the dependency injection logic (which usually happens in cmd/server/main.go).
func InitializeTestHandlers(r *chi.Mux, accountRepo *mocks.MockAccountRepository, userProfileRepo *mocks.MockUserProfileRepository) {
	// We inject the mock repositories into the handlers.
	// Note: We need type casting because the handlers expect the interface,
	// but we are passing the concrete mock struct.
	var accRepoInterface repository.AccountRepository = accountRepo
	var userRepoInterface userprofilerepository.UserProfileRepository = userProfileRepo

	// WireAccountRoutes would be defined in the handler package.
	// For testing, we call it directly.
	WireAccountRoutes(r, accRepoInterface, userRepoInterface)
}
