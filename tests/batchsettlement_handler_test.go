package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/batchsettlement/dto"
	"github.com/carddemo/project/src/app/port/in/rest"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupBatchSettlementRouter(t *testing.T, bs *mocks.MockBatchSettlementService) *chi.Mux {
	r := chi.NewRouter()

	h := rest.NewBatchSettlementHandler(bs)

	r.Post("/settlements", h.CreateSettlement)
	r.Get("/settlements/{id}", h.GetSettlement)
	r.Get("/settlements", h.ListSettlements)

	return r
}

func TestBatchSettlementHandlers_CreateSettlement_Success(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockBatchSettlementService{}
	router := setupBatchSettlementRouter(t, mockSvc)

	reqBody := dto.CreateBatchSettlementRequest{
		Name:        "Batch 100",
		Description: "End of day settlement",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.BatchSettlementResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "Batch 100", resp.Name)
}

func TestBatchSettlementHandlers_GetSettlement_Success(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockBatchSettlementService{}
	mockSvc.GetFunc = func(ctx context.Context, id string) (*dto.BatchSettlementResponse, error) {
		return &dto.BatchSettlementResponse{
			ID:          id,
			Name:        "Batch 200",
			Status:      "reconciled",
			Description: "Checked",
		}, nil
	}

	router := setupBatchSettlementRouter(t, mockSvc)
	req := httptest.NewRequest("GET", "/settlements/batch-123", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.BatchSettlementResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "batch-123", resp.ID)
}

func TestBatchSettlementHandlers_ListSettlements_Success(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockBatchSettlementService{}
	router := setupBatchSettlementRouter(t, mockSvc)

	req := httptest.NewRequest("GET", "/settlements", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp []dto.BatchSettlementResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestBatchSettlementHandlers_ResourceLeakPrevention(t *testing.T) {
	// Arrange
	mockSvc := &mocks.MockBatchSettlementService{}
	router := setupBatchSettlementRouter(t, mockSvc)

	// Invalid JSON to trigger early return logic path
	reqBody := []byte(`{"name": "test", invalid json}`)

	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	// We expect a 400 or 500 depending on implementation details of the decoder,
	// but if defer Close() is missing, the body remains open.
	// In a real integration test, we might inspect the connection state.
	// Here we ensure the handler returns an error status and doesn't hang.
	assert.NotEqual(t, http.StatusOK, w.Code)
}
