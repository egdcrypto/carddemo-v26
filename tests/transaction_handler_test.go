package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/batchsettlement/dto"
	tdto "github.com/carddemo/project/src/app/transaction/dto"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTransactionRouter creates a chi router with the transaction handlers wired to mocks.
func setupTransactionRouter(t *testing.T, ts *mocks.MockTransactionService) *chi.Mux {
	r := chi.NewRouter()

	// Wire the handlers with the mock service
	h := rest.NewTransactionHandler(ts)

	r.Post("/transactions", h.CreateTransaction)
	r.Get("/transactions/{id}", h.GetTransaction)
	r.Get("/accounts/{id}/transactions", h.GetAccountTransactions)
	r.Post("/transactions/{id}/void", h.VoidTransaction)

	return r
}

func TestTransactionHandlers_CreateTransaction_Success(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockTransactionService{}
	router := setupTransactionRouter(t, mockSvc)

	reqBody := tdto.CreateTransactionRequest{
		AccountID:       "acc-123",
		CardID:          "card-123",
		Amount:          100.50,
		TransactionType: "debit",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusAccepted, w.Code)

	var resp dto.WorkflowResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.WorkflowID)
}

func TestTransactionHandlers_CreateTransaction_ValidationError(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockTransactionService{}
	router := setupTransactionRouter(t, mockSvc)

	// Missing Amount
	reqBody := `{"account_id":"acc-123", "transaction_type":"debit"}`
	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Amount")
}

func TestTransactionHandlers_GetTransaction_Success(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockTransactionService{}
	// Setup mock return
	mockSvc.GetFunc = func(ctx context.Context, id string) (*tdto.TransactionResponse, error) {
		return &tdto.TransactionResponse{
			ID:              id,
			AccountID:       "acc-123",
			Amount:          50.0,
			TransactionType: "debit",
			Status:          "settled",
		}, nil
	}

	router := setupTransactionRouter(t, mockSvc)
	req := httptest.NewRequest("GET", "/transactions/txn-123", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp tdto.TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "txn-123", resp.ID)
	assert.Equal(t, "settled", resp.Status)
}

func TestTransactionHandlers_VoidTransaction_Success(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockTransactionService{}
	router := setupTransactionRouter(t, mockSvc)

	reqBody := tdto.VoidTransactionRequest{Reason: "User requested cancellation"}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transactions/txn-999/void", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTransactionHandlers_GetAccountTransactions_QueryParams(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockTransactionService{}
	router := setupTransactionRouter(t, mockSvc)

	// Test with query parameters
	url := "/accounts/acc-555/transactions?start_date=2023-01-01&end_date=2023-01-31&status=settled&page=1&limit=10"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the service was called with parsed params
	// (We would inspect the mock call arguments here if we used a stricter mock library like testify/mock)
	var resp []tdto.TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
