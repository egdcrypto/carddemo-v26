package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
	userprofile_dto "github.com/carddemo/project/src/app/userprofile/dto"
	"github.com/carddemo/project/mocks"
	"github.com/go-chi/chi/v5"
)

// setupRouter creates a test router with injected mock repositories.
func setupRouter(accountRepo repository.AccountRepository, profileRepo userprofile_repository.UserProfileRepository) *chi.Mux {
	r := chi.NewRouter()
	
	// We expect the NewHandler function to exist in the rest package
	handler := rest.NewAccountHandler(accountRepo, profileRepo)
	
	r.Mount("/accounts", handler.Routes())
	return r
}

// TestCreateAccount_Success tests the happy path for account creation.
func TestCreateAccount_Success(t *testing.T) {
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccountRepo, mockProfileRepo)

	reqBody := dto.CreateAccountRequest{
		UserProfileID: "user-123",
		AccountType:   "SAVINGS",
		Status:        "ACTIVE",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

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
	if resp.Status != "ACTIVE" {
		t.Errorf("Expected status ACTIVE, got %s", resp.Status)
	}
}

// TestCreateAccount_ValidationError tests request validation.
func TestCreateAccount_ValidationError(t *testing.T) {
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccountRepo, mockProfileRepo)

	reqBody := `{"user_profile_id": ""}` // Missing required fields

	req := httptest.NewRequest("POST", "/accounts", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestGetAccount_Success tests retrieving an account.
func TestGetAccount_Success(t *testing.T) {
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()

	// Pre-populate mock
	acc := &model.Account{
		ID:            "acc-123",
		UserProfileID: "user-1",
		AccountType:   "CHECKING",
		Status:        "ACTIVE",
		Version:       1,
	}
	mockAccountRepo.Save(acc)

	r := setupRouter(mockAccountRepo, mockProfileRepo)

	req := httptest.NewRequest("GET", "/accounts/acc-123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp dto.AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.ID != "acc-123" {
		t.Errorf("Expected ID acc-123, got %s", resp.ID)
	}
}

// TestGetAccount_NotFound tests 404 response.
func TestGetAccount_NotFound(t *testing.T) {
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(mockAccountRepo, mockProfileRepo)

	req := httptest.NewRequest("GET", "/accounts/non-existent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestUpdateAccountStatus_Success tests updating an account.
func TestUpdateAccountStatus_Success(t *testing.T) {
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()

	// Pre-populate
	acc := &model.Account{
		ID:            "acc-123",
		UserProfileID: "user-1",
		Status:        "ACTIVE",
		Version:       1,
	}
	mockAccountRepo.Save(acc)

	r := setupRouter(mockAccountRepo, mockProfileRepo)

	reqBody := dto.UpdateAccountStatusRequest{NewStatus: "SUSPENDED", Reason: "Security"}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/accounts/acc-123", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp dto.AccountResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Status != "SUSPENDED" {
		t.Errorf("Expected status SUSPENDED, got %s", resp.Status)
	}
}

// TestDeleteAccount_Success tests deleting an account.
func TestDeleteAccount_Success(t *testing.T) {
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()

	acc := &model.Account{ID: "acc-123", Status: "ACTIVE"}
	mockAccountRepo.Save(acc)

	r := setupRouter(mockAccountRepo, mockProfileRepo)

	req := httptest.NewRequest("DELETE", "/accounts/acc-123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	_, err := mockAccountRepo.Get("acc-123")\t
	// The mock repo implementation in setup returns nil, nil on not found.
	// Our Delete method in mock should remove the key.
	if err != nil {
		t.Errorf("Expected no error fetching deleted account, got %v", err)
	}
}

// TestGetUserProfile_Success tests retrieving a user profile via account.
func TestGetUserProfile_Success(t *testing.T) {
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()

	// Pre-populate profile linked to account
	profile := &userprofile_model.UserProfile{
		ID:        "prof-123",
		AccountID: "acc-123",
		Email:     "test@example.com",
	}
	mockProfileRepo.Save(profile)

	r := setupRouter(mockAccountRepo, mockProfileRepo)

	req := httptest.NewRequest("GET", "/accounts/acc-123/profile", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp userprofile_dto.UserProfileResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.AccountID != "acc-123" {
		t.Errorf("Expected AccountID acc-123, got %s", resp.AccountID)
	}
}

// TestUpdateUserProfile_Success tests updating a user profile.
func TestUpdateUserProfile_Success(t *testing.T) {
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockProfileRepo := mocks.NewMockUserProfileRepository()

	profile := &userprofile_model.UserProfile{
		ID:        "prof-123",
		AccountID: "acc-123",
		FirstName: "Old",
		LastName:  "Name",
		Email:     "old@example.com",
	}
	mockProfileRepo.Save(profile)

	r := setupRouter(mockAccountRepo, mockProfileRepo)

	reqBody := userprofile_dto.UpdateUserProfileRequest{
		FirstName: "New",
		LastName:  "Name",
		Email:     "new@example.com",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/accounts/acc-123/profile", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp userprofile_dto.UserProfileResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.FirstName != "New" {
		t.Errorf("Expected FirstName New, got %s", resp.FirstName)
	}
}
