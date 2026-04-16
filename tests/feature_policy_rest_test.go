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
	"github.com/carddemo/project/src/domain/card/model"
	cardpolicyrepo "github.com/carddemo/project/src/domain/cardpolicy/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPolicyRouter(
	policyRepo cardpolicyrepo.CardPolicyRepository,
	cardRepo mocks.MockCardRepository,
) *chi.Mux {
	r := chi.NewRouter()
	rest.RegisterPolicyRoutes(r, policyRepo, cardRepo)
	return r
}

// TestPolicyEndpoints_GetPolicy tests GET /policies/{id}.
func TestPolicyEndpoints_GetPolicy(t *testing.T) {
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()
	mockCardRepo := mocks.NewMockCardRepository()

	router := setupPolicyRouter(mockPolicyRepo, mockCardRepo)

	t.Run("Success: Get Policy by ID", func(t *testing.T) {
		policy := &model.CardPolicy{
			ID:         "pol_001",
			CardID:     "card_001",
			DailyLimit: 10000,
		}
		mockPolicyRepo.Save(policy)

		req := httptest.NewRequest("GET", "/policies/pol_001", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CardPolicyResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "pol_001", resp.ID)
		assert.Equal(t, 10000, resp.DailyLimit)
	})

	t.Run("Failure: Policy Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/policies/ghost", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestPolicyEndpoints_UpdatePolicy tests PUT /policies/{id}.
func TestPolicyEndpoints_UpdatePolicy(t *testing.T) {
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()
	mockCardRepo := mocks.NewMockCardRepository()

	router := setupPolicyRouter(mockPolicyRepo, mockCardRepo)

	t.Run("Success: Update Policy Limits", func(t *testing.T) {
		policy := &model.CardPolicy{
			ID:         "pol_002",
			CardID:     "card_002",
			DailyLimit: 5000,
		}
		mockPolicyRepo.Save(policy)

		reqBody := dto.UpdateCardPolicyRequest{DailyLimit: 8000, WeeklyLimit: 50000}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/policies/pol_002", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp dto.CardPolicyResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 8000, resp.DailyLimit)
		assert.Equal(t, 50000, resp.WeeklyLimit)
	})

	t.Run("Failure: Invalid Payload (Negative Limit)", func(t *testing.T) {
		reqBody := map[string]int{"daily_limit": -100}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/policies/pol_002", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestPolicyEndpoints_GetPoliciesByAccount tests GET /accounts/{id}/policies.
func TestPolicyEndpoints_GetPoliciesByAccount(t *testing.T) {
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()
	mockCardRepo := mocks.NewMockCardRepository()

	router := setupPolicyRouter(mockPolicyRepo, mockCardRepo)

	t.Run("Success: Get Policies for Account", func(t *testing.T) {
		// Seed cards and policies for acc_101
		card1 := &model.Card{ID: "c1", AccountID: "acc_101"}
		card2 := &model.Card{ID: "c2", AccountID: "acc_101"}
		mockCardRepo.Save(card1)
		mockCardRepo.Save(card2)

		pol1 := &model.CardPolicy{ID: "p1", CardID: "c1", DailyLimit: 100}
		pol2 := &model.CardPolicy{ID: "p2", CardID: "c2", DailyLimit: 200}
		mockPolicyRepo.Save(pol1)
		mockPolicyRepo.Save(pol2)

		req := httptest.NewRequest("GET", "/accounts/acc_101/policies", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dto.CardPolicyResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 2)
	})
}
