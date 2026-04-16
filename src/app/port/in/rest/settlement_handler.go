package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/carddemo/project/src/app/batchsettlement/dto"
	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/batchsettlement/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// BatchSettlementHandler handles HTTP requests for Batch Settlements.
type BatchSettlementHandler struct {
	repo     repository.BatchSettlementRepository
	validate *validator.Validate
	logger   *zap.Logger
}

// NewBatchSettlementHandler initializes the handler.
func NewBatchSettlementHandler(repo repository.BatchSettlementRepository) *BatchSettlementHandler {
	return &BatchSettlementHandler{
		repo:     repo,
		validate: validator.New(),
		logger:   zap.NewNop(),
	}
}

// RegisterRoutes mounts the routes.
func (h *BatchSettlementHandler) RegisterRoutes(r chi.Router) {
	r.Post("/settlements", h.CreateSettlement)
	r.Get("/settlements/{id}", h.GetSettlement)
	r.Get("/settlements", h.ListSettlements)
}

func (h *BatchSettlementHandler) CreateSettlement(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSettlementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	batchID := "batch_" + time.Now().Format("20060102150405")

	batch := model.NewBatchSettlement(batchID)

	cmd := command.OpenBatchCmd{
		BatchID:    batchID,
		MerchantID: req.MerchantID,
		Amount:     req.Amount,
	}

	if err := batch.Execute(cmd); err != nil {
		h.logger.Error("Failed to execute batch command", zap.Error(err))
		http.Error(w, "Failed to create settlement", http.StatusBadRequest)
		return
	}

	if err := h.repo.Save(batch); err != nil {
		h.logger.Error("Failed to save settlement", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.mapToResponse(batch))
}

func (h *BatchSettlementHandler) GetSettlement(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	batch, err := h.repo.Get(id)
	if err != nil {
		http.Error(w, "Settlement not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.mapToResponse(batch))
}

func (h *BatchSettlementHandler) ListSettlements(w http.ResponseWriter, r *http.Request) {
	batches, err := h.repo.List()
	if err != nil {
		http.Error(w, "Failed to list settlements", http.StatusInternalServerError)
		return
	}

	var resp []dto.SettlementResponse
	for _, b := range batches {
		resp = append(resp, h.mapToResponse(b))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BatchSettlementHandler) mapToResponse(b *model.BatchSettlement) dto.SettlementResponse {
	return dto.SettlementResponse{
		ID:         b.ID,
		MerchantID: b.MerchantID,
		Amount:     b.Amount,
		Status:     b.Status,
	}
}
