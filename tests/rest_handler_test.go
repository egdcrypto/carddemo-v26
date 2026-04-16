package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accountdto "github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"

	"github.com/carddemo/project/mocks"
)

// setupRouter creates a test router with mocks injected.
func setupRouter(accountRepo repository.AccountRepository, userRepo repository.UserProfileRepository) *chi.Mux {
	r := chi.NewRouter()
	handler := rest.NewAccountHandler(accountRepo, userRepo)
	handler.RegisterRoutes(r)
	return r
}

func TestAccountEndpoints_CreateAccount(t *testing.T) {
	mockAccRepo := mocks.NewMockAccountRepository()
	mockUserRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccRepo, mockUserRepo)

	t.Run("Success 201", func(t *testing.T) {
		body := map[string]string{
			"user_profile_id": "user-123",
			"account_type":    "checking",
			"status":          "active",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "id")
	})

	t.Run("Failure 400 Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Failure 400 Validation Error", func(t *testing.T) {
		body := map[string]string{
			"status": "",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAccountEndpoints_GetByAccountID(t *testing.T) {
	mockAccRepo := mocks.NewMockAccountRepository()
	mockUserRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccRepo, mockUserRepo)

	t.Run("Success 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/12345", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp accountdto.AccountResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "12345", resp.ID)
	})

	t.Run("Failure 404 Not Found (Simulated)", func(t *testing.T) {
		// In the mock handler provided in Red Phase, ID="404" triggers 404
		req := httptest.NewRequest("GET", "/accounts/404", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Failure 500 Internal Error (Simulated)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/500", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAccountEndpoints_UpdateAccount(t *testing.T) {
	mockAccRepo := mocks.NewMockAccountRepository()
	mockUserRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccRepo, mockUserRepo)

	t.Run("Success 200", func(t *testing.T) {
		body := map[string]string{"status": "frozen", "reason": "security"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/123", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Failure 400 Validation", func(t *testing.T) {
		body := map[string]string{"status": ""}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/123", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Failure 404 Not Found", func(t *testing.T) {
		body := map[string]string{"status": "closed"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/404", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAccountEndpoints_DeleteAccount(t *testing.T) {
	mockAccRepo := mocks.NewMockAccountRepository()
	mockUserRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccRepo, mockUserRepo)

	t.Run("Success 204 No Content", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/accounts/123", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Failure 404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/accounts/404", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAccountEndpoints_GetProfileByAccountID(t *testing.T) {
	mockAccRepo := mocks.NewMockAccountRepository()
	mockUserRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccRepo, mockUserRepo)

	t.Run("Success 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/acc-123/profile", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp accountdto.UserProfileResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "profile-acc-123", resp.ID)
	})

	t.Run("Failure 404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/404/profile", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAccountEndpoints_UpdateProfileByAccountID(t *testing.T) {
	mockAccRepo := mocks.NewMockAccountRepository()
	mockUserRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccRepo, mockUserRepo)

	t.Run("Success 200", func(t *testing.T) {
		body := map[string]string{
			"first_name": "John",
			"last_name":  "Doe",
			"email":       "john@example.com",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/acc-123/profile", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Failure 400 Validation", func(t *testing.T) {
		body := map[string]string{
			"first_name": "",
			"last_name":  "",
			"email":       "",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/accounts/acc-123/profile", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
