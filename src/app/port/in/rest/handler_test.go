package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/src/app/account/dto"
	userprofiledto "github.com/carddemo/project/src/app/userprofile/dto"
	"github.com/carddemo/project/src/domain/account/model"
	userprofilemodel "github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/mocks"
	"github.com/go-chi/chi/v5"
)

// setupRouter initializes the chi router with the handlers under test,
// injecting the provided mocks.
func setupRouter(accountRepo *mocks.MockAccountRepository, userProfileRepo *mocks.MockUserProfileRepository) *chi.Mux {
	r := chi.NewRouter()
	// In the real code, handlers would be wired via dependency injection.
	// Here we inject them directly via a constructor helper for testing.
	InitializeTestHandlers(r, accountRepo, userProfileRepo)
	return r
}

func TestPostAccounts(t *testing.T) {
	// Setup
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	// Seed a user profile since account creation links to one
	userProfile := &userprofilemodel.UserProfile{
		ID:    "user-123",
		Email: "test@example.com",
	}
	_ = userProfileRepo.Save(userProfile)

	body := map[string]interface{}{
		"user_profile_id": "user-123",
		"account_type":    "SAVINGS",
		"initial_status":  "ACTIVE",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}

	var response dto.AccountResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccountType != "SAVINGS" {
		t.Errorf("Expected account type SAVINGS, got %s", response.AccountType)
	}
	if response.Status != "ACTIVE" {
		t.Errorf("Expected status ACTIVE, got %s", response.Status)
	}
	if response.ID == "" {
		t.Error("Expected ID to be set")
	}

	// Verify Repository State
	// In a real scenario, we might check length or specific ID, but this validates integration.
}

func TestPostAccounts_ValidationError(t *testing.T) {
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	body := map[string]interface{}{
		"user_profile_id": "", // Invalid
		"account_type":    "INVALID_TYPE",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestGetAccountByID(t *testing.T) {
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	// Seed Data
	account := &model.Account{
		ID:            "acc-999",
		UserProfileID: "user-123",
		Status:        "ACTIVE",
		AccountType:   "CHECKING",
	}
	_ = accountRepo.Save(account)

	req, _ := http.NewRequest("GET", "/accounts/acc-999", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response dto.AccountResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != "acc-999" {
		t.Errorf("Expected ID acc-999, got %s", response.ID)
	}
}

func TestGetAccountByID_NotFound(t *testing.T) {
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	req, _ := http.NewRequest("GET", "/accounts/nonexistent", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

func TestPutAccount(t *testing.T) {
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	// Seed
	account := &model.Account{
		ID:            "acc-999",
		UserProfileID: "user-123",
		Status:        "ACTIVE",
		AccountType:   "CHECKING",
	}
	_ = accountRepo.Save(account)

	body := map[string]interface{}{
		"status": "SUSPENDED",
		"reason": "Suspicious activity",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/accounts/acc-999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response dto.AccountResponse
	json.NewDecoder(rec.Body).Decode(&response)

	if response.Status != "SUSPENDED" {
		t.Errorf("Expected status SUSPENDED, got %s", response.Status)
	}
}

func TestDeleteAccount(t *testing.T) {
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	// Seed
	account := &model.Account{
		ID:            "acc-delete-me",
		UserProfileID: "user-123",
		Status:        "ACTIVE",
		AccountType:   "CHECKING",
	}
	_ = accountRepo.Save(account)

	req, _ := http.NewRequest("DELETE", "/accounts/acc-delete-me", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rec.Code)
	}

	// Verify deleted
	_, err := accountRepo.Get("acc-delete-me")
	if err == nil {
		t.Error("Expected repository to return error for deleted account")
	}
}

// --- User Profile Tests ---

func TestGetUserProfile(t *testing.T) {
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	// Seed profile
	profile := &userprofilemodel.UserProfile{
		ID:        "prof-555",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		AccountID: "acc-123",
	}
	_ = userProfileRepo.Save(profile)

	req, _ := http.NewRequest("GET", "/accounts/acc-123/profile", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response userprofiledto.UserProfileResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.FirstName != "John" {
		t.Errorf("Expected FirstName John, got %s", response.FirstName)
	}
}

func TestPutUserProfile(t *testing.T) {
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	// Seed profile
	profile := &userprofilemodel.UserProfile{
		ID:        "prof-555",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		AccountID: "acc-123",
	}
	_ = userProfileRepo.Save(profile)

	body := map[string]interface{}{
		"first_name": "Jane",
		"last_name":  "Smith",
		"email":      "jane@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/accounts/acc-123/profile", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response userprofiledto.UserProfileResponse
	json.NewDecoder(rec.Body).Decode(&response)

	if response.FirstName != "Jane" {
		t.Errorf("Expected FirstName Jane, got %s", response.FirstName)
	}
}

func TestPutUserProfile_ValidationError(t *testing.T) {
	accountRepo := mocks.NewMockAccountRepository()
	userProfileRepo := mocks.NewMockUserProfileRepository()
	r := setupRouter(accountRepo, userProfileRepo)

	profile := &userprofilemodel.UserProfile{
		ID:        "prof-555",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		AccountID: "acc-123",
	}
	_ = userProfileRepo.Save(profile)

	body := map[string]interface{}{
		"first_name": "J", // Too short
		"email":      "not-an-email",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/accounts/acc-123/profile", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}
