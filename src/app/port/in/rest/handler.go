package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	batchDto "github.com/carddemo/project/src/app/batchsettlement/dto"
	txDto "github.com/carddemo/project/src/app/transaction/dto"

	batchModel "github.com/carddemo/project/src/domain/batchsettlement/model"
	batchRepo "github.com/carddemo/project/src/domain/batchsettlement/repository"

	txnCommand "github.com/carddemo/project/src/domain/transaction/command"
	txnModel "github.com/carddemo/project/src/domain/transaction/model"
	txnRepo "github.com/carddemo/project/src/domain/transaction/repository"

	"github.com/go-chi/chi/v5"
)

// Handler acts as the REST controller for Transactions and Batch Settlements.
type Handler struct {
	txnRepo   txnRepo.TransactionRepository
	batchRepo batchRepo.BatchSettlementRepository
}

// NewHandler creates a new REST handler.
func NewHandler(tr txnRepo.TransactionRepository, br batchRepo.BatchSettlementRepository) *Handler {
	return &Handler{
		txnRepo:   tr,
		batchRepo: br,
	}
}

// RegisterRoutes mounts the handler's routes onto the chi router.
// Middleware can be applied here to specific route groups as required by the AC.
func (h *Handler) RegisterRoutes(r *chi.Mux) {
	// Transaction Routes
	r.Route("/transactions", func(r chi.Router) {
		r.Use(loggingMiddleware) // Example middleware chain
		r.Post("/", h.PostTransaction)
		r.Get("/{txnID}", h.GetTransaction)
		r.Post("/{txnID}/void", h.PostTransactionVoid)
	})

	// Account Transaction Routes
	r.Route("/accounts", func(r chi.Router) {
		r.Use(loggingMiddleware)
		r.Get("/{accountID}/transactions", h.GetAccountTransactions)
	})

	// Batch Settlement Routes
	r.Route("/settlements", func(r chi.Router) {
		r.Use(loggingMiddleware)
		r.Post("/", h.PostBatch)
		r.Get("/{batchID}", h.GetBatch)
		r.Get("/", h.ListBatches)
	})
}

// --- Transaction Handlers ---

// PostTransaction handles the creation of a new transaction.
func (h *Handler) PostTransaction(w http.ResponseWriter, r *http.Request) {
	var req txDto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validation
	if req.AccountID == "" || req.CardID == "" || req.TransactionType == "" {
		http.Error(w, "Missing required fields: account_id, card_id, transaction_type", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	// Map to Domain Command
	cmd := txnCommand.SubmitTransactionCmd{
		TransactionID:   generateID("txn"),
		AccountID:       req.AccountID,
		CardID:          req.CardID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		AccountStatus:   "Active", // In a real app, fetch from Account service
	}

	// Create Aggregate
	txn := txnModel.NewTransaction(cmd.TransactionID, cmd.AccountID, cmd.CardID, cmd.Amount, cmd.TransactionType)

	// Save to Repository
	if err := h.txnRepo.Save(txn); err != nil {
		http.Error(w, "Failed to save transaction", http.StatusInternalServerError)
		return
	}

	// Trigger Temporal Workflow (Async - side effect)
	// In a real app: h.workflowClient.Execute(txn.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(txDto.TransactionResponse{
		ID:              txn.ID,
		AccountID:       txn.AccountID,
		CardID:          txn.CardID,
		Amount:          txn.Amount,
		TransactionType: txn.Type,
		Status:          txn.Status,
		CreatedAt:       txn.CreatedAt.Format(time.RFC3339),
	})
}

// GetTransaction retrieves a single transaction by ID.
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	txnID := chi.URLParam(r, "txnID")
	if txnID == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	txn, err := h.txnRepo.Get(txnID)
	if err != nil || txn == nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txDto.TransactionResponse{
		ID:              txn.ID,
		AccountID:       txn.AccountID,
		CardID:          txn.CardID,
		Amount:          txn.Amount,
		TransactionType: txn.Type,
		Status:          txn.Status,
		CreatedAt:       txn.CreatedAt.Format(time.RFC3339),
	})
}

// GetAccountTransactions retrieves transactions for a specific account.
func (h *Handler) GetAccountTransactions(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	if accountID == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	// Query Params
	_ = r.URL.Query().Get("status")
	// In a real implementation, pass these filters to the repo.

	allTxns, err := h.txnRepo.List()
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	// Filter manually for InMemory repo compliance (CQRS Read Side would handle this in DB)
	var results []*txnModel.Transaction
	for _, t := range allTxns {
		if t.AccountID == accountID {
			results = append(results, t)
		}
	}

	resp := make([]txDto.TransactionResponse, len(results))
	for i, t := range results {
		resp[i] = txDto.TransactionResponse{
			ID:              t.ID,
			AccountID:       t.AccountID,
			CardID:          t.CardID,
			Amount:          t.Amount,
			TransactionType: t.Type,
			Status:          t.Status,
			CreatedAt:       t.CreatedAt.Format(time.RFC3339),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PostTransactionVoid handles voiding (reversing) a transaction.
func (h *Handler) PostTransactionVoid(w http.ResponseWriter, r *http.Request) {
	txnID := chi.URLParam(r, "txnID")
	if txnID == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	var req txDto.VoidTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	txn, err := h.txnRepo.Get(txnID)
	if err != nil || txn == nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Execute Domain Logic
	cmd := txnCommand.ReverseTransactionCmd{
		TransactionID: txnID,
		Reason:        req.Reason,
		Amount:        txn.Amount,
		AccountStatus: "Active",
	}

	if err := txn.Handle(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.txnRepo.Save(txn); err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txDto.TransactionResponse{
		ID:              txn.ID,
		AccountID:       txn.AccountID,
		CardID:          txn.CardID,
		Amount:          txn.Amount,
		TransactionType: txn.Type,
		Status:          txn.Status,
		CreatedAt:       txn.CreatedAt.Format(time.RFC3339),
	})
}

// --- Batch Settlement Handlers ---

// PostBatch handles the creation of a new batch settlement.
func (h *Handler) PostBatch(w http.ResponseWriter, r *http.Request) {
	var req batchDto.CreateBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validation
	if req.MerchantID == "" || req.SettlementDate == "" {
		http.Error(w, "Missing required fields: merchant_id, settlement_date", http.StatusBadRequest)
		return
	}

	settlementDate, err := time.Parse(time.RFC3339, req.SettlementDate)
	if err != nil {
		http.Error(w, "Invalid date format. Use RFC3339 (e.g. 2023-10-27T10:00:00Z)", http.StatusBadRequest)
		return
	}

	// Create Aggregate
	batch := batchModel.NewBatchSettlement(generateID("batch"), req.MerchantID, settlementDate)

	// Save
	if err := h.batchRepo.Save(batch); err != nil {
		http.Error(w, "Failed to create batch", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(batchDto.BatchSettlementResponse{
		ID:        batch.ID,
		Status:    batch.Status,
		CreatedAt: batch.CreatedAt.Format(time.RFC3339),
	})
}

// GetBatch retrieves a single batch settlement by ID.
func (h *Handler) GetBatch(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchID")
	if batchID == "" {
		http.Error(w, "Missing batch ID", http.StatusBadRequest)
		return
	}

	batch, err := h.batchRepo.Get(batchID)
	if err != nil || batch == nil {
		http.Error(w, "Batch not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batchDto.BatchSettlementResponse{
		ID:        batch.ID,
		Status:    batch.Status,
		CreatedAt: batch.CreatedAt.Format(time.RFC3339),
	})
}

// ListBatches retrieves all batch settlements.
func (h *Handler) ListBatches(w http.ResponseWriter, r *http.Request) {
	batches, err := h.batchRepo.List()
	if err != nil {
		http.Error(w, "Failed to retrieve batches", http.StatusInternalServerError)
		return
	}

	resp := make([]batchDto.BatchSettlementResponse, len(batches))
	for i, b := range batches {
		resp[i] = batchDto.BatchSettlementResponse{
			ID:        b.ID,
			Status:    b.Status,
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- Helpers ---

func generateID(prefix string) string {
	return prefix + "_" + strconv.FormatInt(time.Now().UnixNano(), 36)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log logic here
		next.ServeHTTP(w, r)
	})
}

func init() {
	// Ensure batchModel command can be handled (if not already initialized in the model package)
	var _ = batchModel.BatchSettlementAggregate.Handle
}
