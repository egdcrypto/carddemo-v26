package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/app/account/handler"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/shared"
	mocks "github.com/carddemo/project/tests/mocks"
	"github.com/go-chi/chi/v5"
)

// Global setup for router and mocks
var testRouter *chi.Mux
var mockAccountRepo *mocks.MockAccountRepository
var mockUserProfileRepo *mocks.MockUserProfileRepository

func setupTestRouter() {
	testRouter = chi.NewRouter()

	// Initialize Mocks
	mockAccountRepo = mocks.NewMockAccountRepository()
	mockUserProfileRepo = mocks.NewMockUserProfileRepository()

	// We would normally wire the repos into the handler via constructor injection.
	// For the Red phase, we instantiate the handler (which is currently empty).
	// We test against the public HTTP interface.

	h := handler.NewAccountHandler()
	h.RegisterRoutes(testRouter)
}

// Test Account Endpoints
func TestAccountEndpoints(t *testing.T) {
	setupTestRouter()

	t.Run("POST /accounts - 201 Created", func(t *testing.T) {
		body := dto.CreateAccountRequest{
			UserProfileID: "user-123",
			Status:       "ACTIVE",
			AccountType:  "CHECKING",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var resp dto.AccountResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if resp.ID == "" {
			t.Error("Expected ID to be set")
		}
	})

	t.Run("POST /accounts - 400 Bad Request (Validation)", func(t *testing.T) {
		body := dto.CreateAccountRequest{
			// Missing required fields
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}

		var errResp dto.ErrorResponse
		json.NewDecoder(w.Body).Decode(&errResp)
		if errResp.Error == "" {
			t.Error("Expected error message")
		}
	})

	t.Run("GET /accounts/{id} - 200 OK", func(t *testing.T) {
		// Setup mock data
		acct := &model.Account{
			AggregateRoot: shared.AggregateRoot{ID: "acct-1", Version: 1},
			UserProfileID:  "user-1",
			Status:        "ACTIVE",
			AccountType:   "SAVINGS",
		}
		mockAccountRepo.Save(acct)

		req := httptest.NewRequest("GET", "/accounts/acct-1", nil)
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp dto.AccountResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if resp.ID != "acct-1" {
			t.Errorf("Expected ID acct-1, got %s", resp.ID)
		}
	})

	t.Run("GET /accounts/{id} - 404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/non-existent", nil)
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})

	t.Run("PUT /accounts/{id} - 200 OK", func(t *testing.T) {
		body := dto.UpdateAccountRequest{
			Status: "SUSPENDED",
			Reason: "Security check",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/acct-1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("PUT /accounts/{id} - 400 Bad Request", func(t *testing.T) {
		body := dto.UpdateAccountRequest{ Status: "INVALID" }
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/acct-1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("DELETE /accounts/{id} - 204 No Content", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/accounts/acct-1", nil)
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}
	})
}

// Test UserProfile Endpoints
func TestUserProfileEndpoints(t *testing.T) {
	setupTestRouter()

	t.Run("GET /accounts/{id}/profile - 200 OK", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/acct-1/profile", nil)
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("PUT /accounts/{id}/profile - 200 OK", func(t *testing.T) {
		body := dto.UpdateUserProfileRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/acct-1/profile", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("PUT /accounts/{id}/profile - 400 Validation Error", func(t *testing.T) {
		body := dto.UpdateUserProfileRequest{
			Email: "not-an-email",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/acct-1/profile", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}
