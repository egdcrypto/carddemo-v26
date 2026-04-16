package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/card/dto"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/card/repository"
	cardpolicyrepo "github.com/carddemo/project/src/domain/cardpolicy/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRouter initializes the chi router with controllers for testing.
// This represents the main.go wiring.
func setupRouter(
	cardRepo repository.CardRepository,
	accountRepo repository.AccountRepository,
	policyRepo cardpolicyrepo.CardPolicyRepository,
) *chi.Mux {
	r := chi.NewRouter()

	// In a real app, we would use a factory to wire the controllers.
	// For testing, we instantiate the handlers directly, injecting mocks.
	// Assuming we have a RegisterCardRoutes function in the rest package.
	rest.RegisterCardRoutes(r, cardRepo, accountRepo, policyRepo)

	return r
}

// TestCardEndpoints_POST_IssueCard tests the POST /accounts/{id}/cards endpoint.
func TestCardEndpoints_POST_IssueCard(t *testing.T) {
	// Setup Mocks
	mockCardRepo := mocks.NewMockCardRepository()
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()

	// Seed a valid account
	account := &model.Account{
		ID:     "acc_123",
		Status: "active",
	}
	mockAccountRepo.Save(account)

	router := setupRouter(mockCardRepo, mockAccountRepo, mockPolicyRepo)

	t.Run("Success: Issue Card", func(t *testing.T) {
		reqBody := dto.IssueCardRequest{
			CardType: "physical",
			SpendingLimits: map[string]int{
				"daily":  5000,
				"weekly": 20000,
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/accounts/acc_123/cards", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp dto.CardResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "acc_123", resp.AccountID)
		assert.Equal(t, "physical", resp.CardType)
		assert.Equal(t, "active", resp.Status) // Should default to active
		assert.NotEmpty(t, resp.ID)
	})

	t.Run("Failure: Account Not Found", func(t *testing.T) {
		reqBody := dto.IssueCardRequest{CardType: "virtual"}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/accounts/nonexistent/cards", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Failure: Invalid Payload (Bad CardType)", func(t *testing.T) {
		reqBody := map[string]interface{}{"card_type": "plastic"} // Invalid type
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/accounts/acc_123/cards", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var errResp dto.ErrorResponse
		json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "card_type")
	})
}

// TestCardEndpoints_GetCard tests GET /cards/{id}.
func TestCardEndpoints_GetCard(t *testing.T) {
	mockCardRepo := mocks.NewMockCardRepository()
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()

	router := setupRouter(mockCardRepo, mockAccountRepo, mockPolicyRepo)

	t.Run("Success: Get Card by ID", func(t *testing.T) {
		// Seed a card
		card := &model.Card{
			ID:        "card_001",
			AccountID: "acc_123",
			CardType:  "virtual",
			Status:    "active",
		}
		mockCardRepo.Save(card)

		req := httptest.NewRequest("GET", "/cards/card_001", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CardResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "card_001", resp.ID)
		assert.Equal(t, "virtual", resp.CardType)
	})

	t.Run("Failure: Card Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cards/ghost", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestCardEndpoints_UpdateStatus tests PUT /cards/{id}/status.
func TestCardEndpoints_UpdateStatus(t *testing.T) {
	mockCardRepo := mocks.NewMockCardRepository()
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()

	router := setupRouter(mockCardRepo, mockAccountRepo, mockPolicyRepo)

	t.Run("Success: Update Status", func(t *testing.T) {
		card := &model.Card{ID: "card_002", Status: "active"}
		mockCardRepo.Save(card)

		reqBody := dto.UpdateCardStatusRequest{Status: "blocked"}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/cards/card_002/status", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Failure: Invalid Status", func(t *testing.T) {
		reqBody := map[string]interface{}{"status": "stolen"} // Should be lost_stolen
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/cards/card_002/status", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestCardEndpoints_Activate tests POST /cards/{id}/activate.
func TestCardEndpoints_Activate(t *testing.T) {
	mockCardRepo := mocks.NewMockCardRepository()
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()

	router := setupRouter(mockCardRepo, mockAccountRepo, mockPolicyRepo)

	t.Run("Success: Activate Card", func(t *testing.T) {
		// Create an inactive card
		card := &model.Card{ID: "card_003", Status: "inactive"}
		mockCardRepo.Save(card)

		reqBody := dto.ActivateCardRequest{ActivationCode: "123456"}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/cards/card_003/activate", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CardResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "active", resp.Status)
	})

	t.Run("Failure: Wrong Activation Code", func(t *testing.T) {
		card := &model.Card{ID: "card_004", Status: "inactive"}
		mockCardRepo.Save(card)

		reqBody := dto.ActivateCardRequest{ActivationCode: "000000"}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/cards/card_004/activate", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Expecting 422 Unprocessable Entity or 400 Bad Request depending on logic (validation vs domain logic)
		assert.Contains(t, []int{http.StatusUnprocessableEntity, http.StatusBadRequest}, w.Code)
	})
}
