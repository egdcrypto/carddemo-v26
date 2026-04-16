package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/carddemo/project/src/app/batchsettlement/dto"
	"github.com/go-chi/chi/v5"
)

// BatchSettlementService defines the interface for settlement operations.
type BatchSettlementService interface {
	CreateSettlement(ctx context.Context, req dto.CreateBatchSettlementRequest) (string, error)
	GetSettlement(ctx context.Context, id string) (*dto.BatchSettlementResponse, error)
	ListSettlements(ctx context.Context) ([]*dto.BatchSettlementResponse, error)
}

// BatchSettlementHandler handles HTTP requests for Batch Settlements.
type BatchSettlementHandler struct {
	Service BatchSettlementService
}

// NewBatchSettlementHandler creates a new handler for settlements.
func NewBatchSettlementHandler(svc BatchSettlementService) *BatchSettlementHandler {
	return &BatchSettlementHandler{Service: svc}
}

func (h *BatchSettlementHandler) CreateSettlement(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBatchSettlementRequest
	// Fix: Move defer close to top to prevent resource leak on decode error
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateSettlement(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.BatchSettlementResponse{
		ID:   id,
		Name: req.Name,
	})
}

func (h *BatchSettlementHandler) GetSettlement(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	settlement, err := h.Service.GetSettlement(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Settlement not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settlement)
}

func (h *BatchSettlementHandler) ListSettlements(w http.ResponseWriter, r *http.Request) {
	settlements, err := h.Service.ListSettlements(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settlements)
}
