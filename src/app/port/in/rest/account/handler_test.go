package account_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/app/port/in/rest/account"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRouter initializes the chi router with the handler and a mock repository.
// We use a dependency injection constructor pattern for the handler.
func setupRouter(repo repository.AccountRepository) http.Handler {
	r := chi.NewRouter()
	handler := account.NewAccountHandler(repo)
	handler.RegisterRoutes(r)
	return r
}

func TestAccountHandlers_CreateAccount(t *testing.T) {
	// "Acceptance Criteria: POST /accounts"
	// "Acceptance Criteria: Successful responses return 201"
	// "Acceptance Criteria: Request validation returns 400"

	t.Run("Success: Creates account and returns 201", func(t *testing.T) {
		mockRepo := mocks.NewMockAccountRepository()
		router := setupRouter(mockRepo)

		payload := dto.CreateAccountRequest{
			UserProfileID: "profile-123",
			Status:        "Active",
			AccountType:   "Checking",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response dto.AccountResponse
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, "Active", response.Status)
		assert.Equal(t, "profile-123", response.UserProfileID)
	})

	t.Run("Failure: Returns 400 for invalid JSON", func(t *testing.T) {
		mockRepo := mocks.NewMockAccountRepository()
		router := setupRouter(mockRepo)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewReader([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Failure: Returns 400 for validation errors (missing required field)", func(t *testing.T) {
		mockRepo := mocks.NewMockAccountRepository()
		router := setupRouter(mockRepo)

		payload := dto.CreateAccountRequest{
			UserProfileID: "", // Invalid
			Status:        "Active",
			AccountType:   "Checking",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestAccountHandlers_GetAccount(t *testing.T) {
	// "Acceptance Criteria: GET /accounts/{id}"
	// "Acceptance Criteria: Successful responses return 200"
	// "Acceptance Criteria: Error responses return 404"

	t.Run("Success: Returns account by ID", func(t *testing.T) {
		mockRepo := mocks.NewMockAccountRepository()
		// Seed mock data directly to simulate DB state
		// Note: In a real app, the aggregate would be loaded, here we mock the repo.Get
		// However, since we are in red phase, we define the behavior we expect.
		// We will manually set up the map in the mock for this scenario if we were unit testing,
		// but here we test the handler. We need the mock to return the object.
		
		// Since we don't have the logic to save yet, we can't POST then GET.
		// We will rely on the mock returning the data if the ID matches a specific string, or 404.
		// For the purpose of this test file, we check that the handler calls the repo correctly.
		// However, to make the test robust, let's assume we can 'inject' a state or use a fixed ID.
		// Ideally, we POST first.
		
		router := setupRouter(mockRepo)

		// Let's assume we Create first to populate the mock (since MockRepo is in-memory)
		createPayload := dto.CreateAccountRequest{
			UserProfileID: "user-1",
			Status:        "Active",
			AccountType:   "Savings",
		}
		body, _ := json.Marshal(createPayload)
		createReq := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))
		createResp := httptest.NewRecorder()
		router.ServeHTTP(createResp, createReq) // This populates the mockRepo if implemented correctly
		
		// Extract ID from create response to use in GET
		var createRespDto dto.AccountResponse
		json.Unmarshal(createResp.Body.Bytes(), &createRespDto)
		id := createRespDto.ID

		// Now GET the account
		req := httptest.NewRequest("GET", "/accounts/"+id, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var getRespDto dto.AccountResponse
		err := json.Unmarshal(resp.Body.Bytes(), &getRespDto)
		require.NoError(t, err)
		assert.Equal(t, id, getRespDto.ID)
		assert.Equal(t, "user-1", getRespDto.UserProfileID)
	})

	t.Run("Failure: Returns 404 for non-existent ID", func(t *testing.T) {
		mockRepo := mocks.NewMockAccountRepository()
		router := setupRouter(mockRepo)

		req := httptest.NewRequest("GET", "/accounts/non-existent", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}

func TestAccountHandlers_UpdateAccount(t *testing.T) {
	// "Acceptance Criteria: PUT /accounts/{id}"
	// "Acceptance Criteria: Successful responses return 200"

	t.Run("Success: Updates account status", func(t *testing.T) {
		mockRepo := mocks.NewMockAccountRepository()
		router := setupRouter(mockRepo)

		// 1. Create an account first
		createPayload := dto.CreateAccountRequest{
			UserProfileID: "user-update",
			Status:        "Active",
			AccountType:   "Checking",
		}
		cBody, _ := json.Marshal(createPayload)
		creq := httptest.NewRequest("POST", "/accounts", bytes.NewReader(cBody))
		creq.Header.Set("Content-Type", "application/json")
		crec := httptest.NewRecorder()
		router.ServeHTTP(crec, creq)

		var cResp dto.AccountResponse
		json.Unmarshal(crec.Body.Bytes(), &cResp)
		id := cResp.ID

		// 2. Update it
		updatePayload := dto.UpdateAccountStatusRequest{
			NewStatus: "Suspended",
			Reason:    "Suspicious activity",
		}
		uBody, _ := json.Marshal(updatePayload)
		ureq := httptest.NewRequest("PUT", "/accounts/"+id, bytes.NewReader(uBody))
		ureq.Header.Set("Content-Type", "application/json")
		urec := httptest.NewRecorder()

		router.ServeHTTP(urec, ureq)

		assert.Equal(t, http.StatusOK, urec.Code)
		
		// Verify response body contains new status
		var uResp dto.AccountResponse
		json.Unmarshal(urec.Body.Bytes(), &uResp)
		assert.Equal(t, "Suspended", uResp.Status)
	})
}

func TestAccountHandlers_DeleteAccount(t *testing.T) {
	// "Acceptance Criteria: DELETE /accounts/{id}"
	// "Acceptance Criteria: Successful responses return 204"

	t.Run("Success: Deletes account", func(t *testing.T) {
		mockRepo := mocks.NewMockAccountRepository()
		router := setupRouter(mockRepo)

		// 1. Create an account
		createPayload := dto.CreateAccountRequest{
			UserProfileID: "user-delete",
			Status:        "Active",
			AccountType:   "Checking",
		}
		cBody, _ := json.Marshal(createPayload)
		creq := httptest.NewRequest("POST", "/accounts", bytes.NewReader(cBody))
		creq.Header.Set("Content-Type", "application/json")
		crec := httptest.NewRecorder()
		router.ServeHTTP(crec, creq)

		var cResp dto.AccountResponse
		json.Unmarshal(crec.Body.Bytes(), &cResp)
		id := cResp.ID

		// 2. Delete it
		dreq := httptest.NewRequest("DELETE", "/accounts/"+id, nil)
		drec := httptest.NewRecorder()

		router.ServeHTTP(dreq, drec)

		assert.Equal(t, http.StatusNoContent, drec.Code)
		assert.Equal(t, 0, drec.Body.Len())

		// 3. Verify it's gone
		greq := httptest.NewRequest("GET", "/accounts/"+id, nil)
		grec := httptest.NewRecorder()
		router.ServeHTTP(greq, grec)
		assert.Equal(t, http.StatusNotFound, grec.Code)
	})
}
