package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/userprofile/model"
	userprofile_repo "github.com/carddemo/project/src/domain/userprofile/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// We need to import the handlers once they are implemented.
// For the purpose of this Red Phase test file, we assume a package structure exists.
// If src/app/port/in/rest does not exist, this will fail to compile, which is expected in Red Phase.
// However, to ensure we can write the test NOW, we will check if we can reference the handler.
// If the package is missing, the compiler error serves as the "Failure".

func TestAccountHandlers_RedPhase(t *testing.T) {
	// 1. Setup Mocks
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()

	// 2. Setup Router (Simulating wiring in cmd/server/main.go)
	r := chi.NewRouter()

	// This section will fail to compile until the handlers package is created.
	// Un-commenting the lines below would trigger the compilation error indicating missing packages.
	/*
		handlers.RegisterAccountRoutes(r, mockAccountRepo, mockProfileRepo)
	*/

	// --- TEST CASES ---

	t.Run("POST /accounts - Returns 201 on valid creation", func(t *testing.T) {
		// Setup Domain State (Mocking the "After" state)
		// In a real flow, handler -> service -> aggregate. Here we mock the repo Get return.
		expectedID := "acc-123"
		mockAccountRepo.Save(&model.Account{
			ID:            expectedID,
			UserProfileID: "user-1",
			Status:        "active",
			AccountType:   "savings",
			Version:       1,
		})

		reqBody := dto.CreateAccountRequest{
			UserProfileID: "user-1",
			AccountType:   "savings",
			Status:        "active",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// This will fail because `RegisterAccountRoutes` doesn't exist yet.
		// r.ServeHTTP(w, req)

		// Explicit failure for Red Phase
		t.Fatal("Handler registration is missing or implementation is missing")

		// Future Assertions (Green Phase goals):
		// assert.Equal(t, http.StatusCreated, w.Code)
		// var resp dto.AccountResponse
		// json.NewDecoder(w.Body).Decode(&resp)
		// assert.Equal(t, expectedID, resp.ID)
	})

	t.Run("POST /accounts - Returns 400 on invalid input", func(t *testing.T) {
		reqBody := `{"user_profile_id": ""}` // Missing required fields
		req := httptest.NewRequest("POST", "/accounts", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		t.Fatal("Validation middleware or handler logic is missing")
		// assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GET /accounts/{id} - Returns 200 when found", func(t *testing.T) {
		mockAccountRepo.Save(&model.Account{ID: "acc-found", Status: "active"})

		req := httptest.NewRequest("GET", "/accounts/acc-found", nil)
		w := httptest.NewRecorder()

		t.Fatal("GET handler is missing")
		// r.ServeHTTP(w, req)
		// assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GET /accounts/{id} - Returns 404 when not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/acc-notfound", nil)
		w := httptest.NewRecorder()

		t.Fatal("GET handler is missing")
		// assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DELETE /accounts/{id} - Returns 204 No Content", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/accounts/acc-delete", nil)
		w := httptest.NewRecorder()

		t.Fatal("DELETE handler is missing")
		// assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("PUT /accounts/{id}/profile - Returns 200 when updating profile", func(t *testing.T) {
		// Setup Profile
		mockProfileRepo.Save(&model.UserProfile{ID: "prof-1", AccountID: "acc-1"})

		reqBody := dto.LinkUserToAccountRequest{
			FirstName: "John",
			LastName:  "Doe",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/accounts/acc-1/profile", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		t.Fatal("PUT Profile handler is missing")
		// assert.Equal(t, http.StatusOK, w.Code)
	})
}
