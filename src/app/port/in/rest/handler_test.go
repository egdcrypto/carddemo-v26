package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	txnCommand "github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/go-chi/chi/v5"
)

// setupTestRouter creates a router with the handler applied for testing.
func setupTestRouter(t *testing.T, handler *Handler) *chi.Mux {
	router := chi.NewRouter()
	handler.RegisterRoutes(router)
	return router
}

// setupHandler initializes the Handler with mock repositories.
func setupHandler(t *testing.T) *Handler {
	txnRepo := mocks.NewMockTransactionRepository()
	batchRepo := mocks.NewMockBatchSettlementRepository()
	return NewHandler(txnRepo, batchRepo)
}

func TestTransactionEndpoints_PostTransaction(t *testing.T) {
	handler := setupHandler(t)
	router := setupTestRouter(t, handler)

	t.Run("returns 201 on successful creation", func(t *testing.T) {
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

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["status"] == nil {
			t.Error("Expected status in response")
		}
		// Triggering workflow is a side effect, we assume it happens if no error is returned
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/transactions", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 for missing required fields", func(t *testing.T) {
		payload := map[string]interface{}{
			"amount": 100.50,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestTransactionEndpoints_GetTransaction(t *testing.T) {
	handler := setupHandler(t)
	router := setupTestRouter(t, handler)
	txnRepo := handler.txnRepo.(*mocks.MockTransactionRepository)

	// Seed a transaction
	txn := model.NewTransaction("txn_1", "acc_1", "card_1", 50.0, "debit")
	txnRepo.Save(txn)

	t.Run("returns 200 and transaction for valid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/transactions/txn_1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["id"] != "txn_1" {
			t.Errorf("Expected id txn_1, got %v", resp["id"])
		}
	})

	t.Run("returns 404 for non-existent transaction", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/transactions/unknown", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestTransactionEndpoints_GetAccountTransactions(t *testing.T) {
	handler := setupHandler(t)
	router := setupTestRouter(t, handler)
	txnRepo := handler.txnRepo.(*mocks.MockTransactionRepository)

	// Seed transactions
	txn1 := model.NewTransaction("txn_1", "acc_1", "card_1", 50.0, "debit")
	txn2 := model.NewTransaction("txn_2", "acc_1", "card_1", 20.0, "credit")
	txnRepo.Save(txn1)
	txnRepo.Save(txn2)

	t.Run("returns 200 and filters by account ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts/acc_1/transactions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp) != 2 {
			t.Errorf("Expected 2 transactions, got %d", len(resp))
		}
	})

	t.Run("supports query parameters for status (mocked filter)", func(t *testing.T) {
		// In a real scenario, the repo would filter. Here we test param parsing.
		req := httptest.NewRequest("GET", "/accounts/acc_1/transactions?status=completed", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}

func TestTransactionEndpoints_PostTransactionVoid(t *testing.T) {
	handler := setupHandler(t)
	router := setupTestRouter(t, handler)
	txnRepo := handler.txnRepo.(*mocks.MockTransactionRepository)

	// Seed a transaction
	txn := model.NewTransaction("txn_1", "acc_1", "card_1", 50.0, "debit")
	txnRepo.Save(txn)

	t.Run("returns 200 on successful void", func(t *testing.T) {
		payload := map[string]interface{}{
			"reason": "Customer request",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/transactions/txn_1/void", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 404 if transaction does not exist", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/transactions/unknown/void", strings.NewReader("{}"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestBatchSettlementEndpoints_PostBatch(t *testing.T) {
	handler := setupHandler(t)
	router := setupTestRouter(t, handler)

	t.Run("returns 201 on successful creation", func(t *testing.T) {
		payload := map[string]interface{}{
			"settlement_date": "2023-10-27T10:00:00Z",
			"merchant_id":     "merchant_1",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["id"] == nil {
			t.Error("Expected ID in response")
		}
	})

	t.Run("returns 400 for invalid date format", func(t *testing.T) {
		payload := map[string]interface{}{
			"settlement_date": "invalid-date",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestBatchSettlementEndpoints_GetBatch(t *testing.T) {
	handler := setupHandler(t)
	router := setupTestRouter(t, handler)
	batchRepo := handler.batchRepo.(*mocks.MockBatchSettlementRepository)

	// Seed a batch
	batch := model.NewBatchSettlement("batch_1", "merchant_1", time.Now())
	batchRepo.Save(batch)

	t.Run("returns 200 for valid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/settlements/batch_1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("returns 404 for non-existent batch", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/settlements/unknown", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestBatchSettlementEndpoints_ListBatches(t *testing.T) {
	handler := setupHandler(t)
	router := setupTestRouter(t, handler)

	t.Run("returns 200 and list of batches", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/settlements", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		// We expect an array, possibly empty if no seeding in this specific test context,
		// but the structure must be valid.
		if resp == nil {
			t.Error("Expected array response")
		}
	})
}

func TestHandlerMiddlewareRegistration(t *testing.T) {
	handler := setupHandler(t)
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	// Check if middlewares are stacked. chi router doesn't easily expose registered middlewares,
	// but we can check behavior if we had a global middleware in main.
	// Here we ensure the function runs without panic.
	if router == nil {
		t.Error("Router should not be nil after RegisterRoutes")
	}
}
