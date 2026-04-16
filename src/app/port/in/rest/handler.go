package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/carddemo/project/src/app/batchsettlement/dto"
	batchService "github.com/carddemo/project/src/app/batchsettlement/service"
	"github.com/carddemo/project/src/app/shared"
	"github.com/carddemo/project/src/app/transaction/dto"
	txService "github.com/carddemo/project/src/app/transaction/service"
	"github.com/go-chi/chi/v5"
)

// TransactionHandler handles HTTP requests for Transactions.
type TransactionHandler struct {
	service *txService.TransactionApplication
}

// NewTransactionHandler creates a new REST handler for transactions.
func NewTransactionHandler(service *txService.TransactionApplication) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// RegisterRoutes registers transaction routes on the router.
func (h *TransactionHandler) RegisterRoutes(r chi.Router) {
	r.Post("/transactions", h.Create)
	r.Get("/transactions/{id}", h.Get)
	r.Post("/transactions/{id}/void", h.Reverse)
	r.Get("/accounts/{id}/transactions", h.ListByAccount)
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := shared.ValidateStruct(req); err != nil {
		respondWithValidationErr(w, err)
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing transaction id", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Get(id)
	if err != nil {
		// In a real app, check if NotFound error
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *TransactionHandler) Reverse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing transaction id", http.StatusBadRequest)
		return
	}

	var req dto.ReverseTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := shared.ValidateStruct(req); err != nil {
		respondWithValidationErr(w, err)
		return
	}

	resp, err := h.service.Reverse(id, req)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *TransactionHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		http.Error(w, "missing account id", http.StatusBadRequest)
		return
	}

	// Query params are parsed but repository mock filters simply by account ID
	// In a real implementation, these would be passed to the repository method.
	_ = r.URL.Query().Get("status")
	_ = r.URL.Query().Get("page")
	_ = r.URL.Query().Get("limit")

	resp, err := h.service.ListByAccount(accountID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// BatchSettlementHandler handles HTTP requests for Batch Settlements.
type BatchSettlementHandler struct {
	service *batchService.SettlementApplication
}

// NewBatchSettlementHandler creates a new REST handler for settlements.
func NewBatchSettlementHandler(service *batchService.SettlementApplication) *BatchSettlementHandler {
	return &BatchSettlementHandler{service: service}
}

// RegisterRoutes registers settlement routes on the router.
func (h *BatchSettlementHandler) RegisterRoutes(r chi.Router) {
	r.Post("/settlements", h.Create)
	r.Get("/settlements/{id}", h.Get)
	r.Get("/settlements", h.List)
}

func (h *BatchSettlementHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSettlementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := shared.ValidateStruct(req); err != nil {
		respondWithValidationErr(w, err)
		return
	}

	resp, err := h.service.Create(req)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (h *BatchSettlementHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing settlement id", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Get(id)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *BatchSettlementHandler) List(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.List()
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// --- Helper Functions ---

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	// For simplicity in green phase, returning 500. In real app, check error types.
	respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func respondWithValidationErr(w http.ResponseWriter, err error) {
	msg := err.Error()
	if strings.Contains(msg, "Amount") {
		msg = "Amount must be greater than 0"
	}
	if strings.Contains(msg, "TransactionType") {
		msg = "Transaction type must be debit or credit"
	}
	if strings.Contains(msg, "Currency") {
		msg = "Currency must be 3 characters"
	}
	respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}
