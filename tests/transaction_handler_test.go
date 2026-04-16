package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/app/transaction/dto"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTransactionRouter() (*chi.Mux, *mocks.MockTransactionRepository) {
	repo := mocks.NewMockTransactionRepository()
	handler := rest.NewTransactionHandler(repo)
	router := chi.NewRouter()
	handler.RegisterRoutes(router)
	return router, repo
}

func TestTransactionHandler_CreateTransaction_Success(t *testing.T) {
	router, _ := setupTransactionRouter()

	payload := map[string]interface{}{
		"account_id":       "acc_123",
		"card_id":          "card_123",
		"amount":           100.50,
		"transaction_type": "debit",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// RED PHASE: We expect 201, but currently likely getting 404 or 500 or 200 (empty)
	// Changing expectation to fail against the empty implementation
	assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

	var resp dto.TransactionResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
}

func TestTransactionHandler_CreateTransaction_ValidationError(t *testing.T) {
	router, _ := setupTransactionRouter()

	// Missing amount and invalid type
	payload := map[string]interface{}{
		"account_id":       "acc_123",
		"card_id":          "card_123",
		"amount":           -10,
		"transaction_type": "invalid",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_GetTransaction_NotFound(t *testing.T) {
	router, repo := setupTransactionRouter()

	// Ensure repo is empty
	repo.Delete("does_not_exist")

	req := httptest.NewRequest("GET", "/transactions/does_not_exist", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expecting 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTransactionHandler_GetTransaction_Found(t *testing.T) {
	router, repo := setupTransactionRouter()

	// Setup Mock Data
	// We need to manually build a Transaction aggregate for the mock repo
	// Importing the model is allowed for test setup
	txn := &model.Transaction{
		AggregateBase: shared.AggregateBase{ID: "txn_456"},
		AccountID:     "acc_123",
		CardID:        "card_123",
		Amount:        50.0,
		Type:          "debit",
	}
	repo.Save(txn)

	req := httptest.NewRequest("GET", "/transactions/txn_456", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.TransactionResponse
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "txn_456", resp.ID)
}

func TestTransactionHandler_ListAccountTransactions_QueryParams(t *testing.T) {
	router, _ := setupTransactionRouter()

	req := httptest.NewRequest("GET", "/accounts/acc_123/transactions?status=cleared&page=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
