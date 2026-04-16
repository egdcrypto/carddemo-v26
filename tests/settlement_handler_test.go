package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/carddemo/project/src/app/batchsettlement/dto"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupSettlementRouter() (*chi.Mux, *mocks.MockBatchSettlementRepository) {
	repo := mocks.NewMockBatchSettlementRepository()
	handler := rest.NewBatchSettlementHandler(repo)
	router := chi.NewRouter()
	handler.RegisterRoutes(router)
	return router, repo
}

func TestBatchSettlementHandler_CreateSettlement_Success(t *testing.T) {
	router, _ := setupSettlementRouter()

	payload := map[string]interface{}{
		"merchant_id": "merch_1",
		"amount":      5000.00,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// RED PHASE: Expect implementation to return 201
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.SettlementResponse
	json.NewDecoder(w.Body).Decode(&resp)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "merch_1", resp.MerchantID)
}

func TestBatchSettlementHandler_ListSettlements(t *testing.T) {
	router, repo := setupSettlementRouter()

	// Seed mock data
	batch := &model.BatchSettlement{
		AggregateBase: shared.AggregateBase{ID: "batch_1"},
		MerchantID:    "merch_1",
		Amount:        100.0,
	}
	repo.Save(batch)

	req := httptest.NewRequest("GET", "/settlements", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
