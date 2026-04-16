package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/card/dto"
	"github.com/carddemo/project/src/app/cardpolicy/dto"
	cardrest "github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/domain/card/command"
	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/cardpolicy/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCardEndpoints(t *testing.T) {
	// Setup Mocks
	mockCardRepo := mocks.NewMockCardRepository()
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()
	mockAccountRepo := mocks.NewMockAccountRepository()

	// Seed Data
	ctx := context.Background()
	policyAgg := model.NewCardPolicy("policy-1", "acc-1")
	policyAgg.Handle(command.UpdateCardLimitsCmd{DailyLimit: 1000, SingleTxnLimit: 100})
	mockPolicyRepo.Save(ctx, policyAgg.GetAggregate())

	mockAccountRepo.data["acc-1"] = struct{}{}

	// Setup Router
	deps := &cardrest.Dependencies{
		CardRepo:       mockCardRepo,
		CardPolicyRepo: mockPolicyRepo,
		AccountRepo:    mockAccountRepo,
	}
	router := cardrest.NewHandler(deps)

	t.Run("POST /accounts/{id}/cards - Success (201)", func(t *testing.T) {
		payload := dto.IssueCardRequest{
			AccountID: "acc-1",
			CardType:  "Virtual",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/accounts/acc-1/cards", bytes.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp dto.CardResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "acc-1", resp.AccountID)
		assert.Equal(t, "Virtual", resp.CardType)
	})

	t.Run("POST /accounts/{id}/cards - Validation Fail (400)", func(t *testing.T) {
		payload := dto.IssueCardRequest{
			AccountID: "acc-1",
			CardType:  "InvalidType",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/accounts/acc-1/cards", bytes.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /accounts/{id}/cards - Account Not Found (404)", func(t *testing.T) {
		payload := dto.IssueCardRequest{
			AccountID: "acc-does-not-exist",
			CardType:  "Virtual",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/accounts/acc-random/cards", bytes.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GET /cards/{id} - Success (200)", func(t *testing.T) {
		// Seed a card
		cardAgg := model.IssueCard(command.IssueCardCmd{AccountID: "acc-1", CardType: "Physical"})
		mockCardRepo.Save(ctx, cardAgg)

		req := httptest.NewRequest("GET", "/cards/"+cardAgg.ID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CardResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, cardAgg.ID, resp.ID)
	})

	t.Run("PUT /cards/{id}/status - Success (200)", func(t *testing.T) {
		cardAgg := model.IssueCard(command.IssueCardCmd{AccountID: "acc-1", CardType: "Physical"})
		mockCardRepo.Save(ctx, cardAgg)

		payload := dto.UpdateCardStatusRequest{Status: "Blocked"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("PUT", "/cards/"+cardAgg.ID+"/status", bytes.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CardResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "Blocked", resp.Status)
	})

	t.Run("POST /cards/{id}/activate - Success (200)", func(t *testing.T) {
		cardAgg := model.IssueCard(command.IssueCardCmd{AccountID: "acc-1", CardType: "Virtual"})
		mockCardRepo.Save(ctx, cardAgg)

		payload := dto.ActivateCardRequest{PIN: "1234"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/cards/"+cardAgg.ID+"/activate", bytes.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CardResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "Active", resp.Status)
	})
}

func TestCardPolicyEndpoints(t *testing.T) {
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()
	mockAccountRepo := mocks.NewMockAccountRepository()

	ctx := context.Background()
	policyAgg := model.NewCardPolicy("pol-1", "acc-1")
	mockPolicyRepo.Save(ctx, policyAgg.GetAggregate())
	mockAccountRepo.data["acc-1"] = struct{}{}

	deps := &cardrest.Dependencies{
		CardPolicyRepo: mockPolicyRepo,
		AccountRepo:    mockAccountRepo,
	}
	router := cardrest.NewHandler(deps)

	t.Run("GET /policies/{id} - Success (200)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/policies/pol-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CardPolicyResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, "pol-1", resp.ID)
	})

	t.Run("PUT /policies/{id} - Success (200)", func(t *testing.T) {
		newLimit := 5000
		payload := dto.UpdateCardPolicyRequest{DailyLimit: &newLimit}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("PUT", "/policies/pol-1", bytes.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CardPolicyResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.Equal(t, 5000, resp.DailyLimit)
	})

	t.Run("GET /accounts/{id}/policies - Success (200)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/acc-1/policies", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dto.CardPolicyResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
		assert.GreaterOrEqual(t, len(resp), 1)
	})
}
