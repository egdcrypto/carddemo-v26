package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/carddemo/project/src/app/transaction/dto"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/carddemo/project/src/domain/transaction/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// TransactionHandler handles HTTP requests for Transactions.
type TransactionHandler struct {
	repo     repository.TransactionRepository
	validate *validator.Validate
	logger   *zap.Logger
}

// NewTransactionHandler initializes the handler.
func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		repo:     repo,
		validate: validator.New(),
		logger:   zap.NewNop(), // Replace with real logger in main.go wiring
	}
}

// RegisterRoutes mounts the routes to the chi router.
func (h *TransactionHandler) RegisterRoutes(r chi.Router) {
	r.Post("/transactions", h.CreateTransaction)
	r.Get("/transactions/{id}", h.GetTransaction)
	r.Post("/transactions/{id}/void", h.VoidTransaction)
	r.Get("/accounts/{id}/transactions", h.ListAccountTransactions)
}

// CreateTransaction handles POST /transactions
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In a real app, ID generation is handled by domain or a service (UUID).
	// For this test, we generate a simple string ID.
	txnID := "txn_" + time.Now().Format("20060102150405")

	cmd := command.SubmitTransactionCmd{
		TransactionID:   txnID,
		AccountID:       req.AccountID,
		CardID:          req.CardID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		AccountStatus:   "Active", // Simplified: assuming active for success path
	}

	// Create Aggregate
	txn := model.NewTransaction(txnID)

	// Execute Command
	if err := txn.Execute(cmd); err != nil {
		h.logger.Error("Failed to execute transaction command", zap.Error(err))
		http.Error(w, "Failed to process transaction", http.StatusBadRequest)
		return
	}

	// Persist
	if err := h.repo.Save(txn); err != nil {
		h.logger.Error("Failed to save transaction", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.mapToResponse(txn))
}

// GetTransaction handles GET /transactions/{id}
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	txn, err := h.repo.Get(id)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.mapToResponse(txn))
}

// VoidTransaction handles POST /transactions/{id}/void
func (h *TransactionHandler) VoidTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.VoidTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	txn, err := h.repo.Get(id)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	cmd := command.ReverseTransactionCmd{
		TransactionID: id,
		Reason:        req.Reason,
		Amount:        txn.Amount, // Preserve amount
	}

	if err := txn.Execute(cmd); err != nil {
		http.Error(w, "Failed to void transaction", http.StatusBadRequest)
		return
	}

	if err := h.repo.Save(txn); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.mapToResponse(txn))
}

// ListAccountTransactions handles GET /accounts/{id}/transactions
func (h *TransactionHandler) ListAccountTransactions(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	// Note: The mock repository List() returns all, so we filter in memory for the mock.
	// In a real SQL/Mongo scenario, this query would be pushed to the repo layer.

	allTxns, err := h.repo.List()
	if err != nil {
		http.Error(w, "Failed to list transactions", http.StatusInternalServerError)
		return
	}

	var result []*model.Transaction
	for _, t := range allTxns {
		if t.AccountID == accountID {
			result = append(result, t)
		}
	}

	var resp []dto.TransactionResponse
	for _, t := range result {
		resp = append(resp, h.mapToResponse(t))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *TransactionHandler) mapToResponse(t *model.Transaction) dto.TransactionResponse {
	return dto.TransactionResponse{
		ID:              t.ID,
		AccountID:       t.AccountID,
		CardID:          t.CardID,
		Amount:          t.Amount,
		TransactionType: t.Type,
		Status:          t.Status,
	}
}
