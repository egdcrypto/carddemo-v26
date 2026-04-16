package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/carddemo/project/src/app/shared/dto"
	tdto "github.com/carddemo/project/src/app/transaction/dto"
	"github.com/go-chi/chi/v5"
)

// TransactionService defines the interface for transaction operations.
type TransactionService interface {
	CreateTransaction(ctx context.Context, req tdto.CreateTransactionRequest) (string, error)
	GetTransaction(ctx context.Context, id string) (*tdto.TransactionResponse, error)
	ListAccountTransactions(ctx context.Context, accountID string, params dto.QueryParams) ([]*tdto.TransactionResponse, error)
	VoidTransaction(ctx context.Context, id string, reason string) error
}

// TransactionHandler handles HTTP requests for Transactions.
type TransactionHandler struct {
	Service TransactionService
}

// NewTransactionHandler creates a new handler for transactions.
func NewTransactionHandler(svc TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: svc}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req tdto.CreateTransactionRequest
	// Fix: Move defer close to top to prevent resource leak on decode error
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic Validation (reflecting test requirements for 'Amount')
	if req.Amount == 0 {
		http.Error(w, "Field validation error: Amount is required", http.StatusBadRequest)
		return
	}

	workflowID, err := h.Service.CreateTransaction(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(tdto.WorkflowResponse{WorkflowID: workflowID})
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	txn, err := h.Service.GetTransaction(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txn)
}

func (h *TransactionHandler) GetAccountTransactions(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	params := parseQueryParams(r.URL.Query())

	txns, err := h.Service.ListAccountTransactions(r.Context(), accountID, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txns)
}

func (h *TransactionHandler) VoidTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req tdto.VoidTransactionRequest
	// Fix: Move defer close to top
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.Service.VoidTransaction(r.Context(), id, req.Reason); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// parseQueryParams extracts query parameters from the URL query string.
func parseQueryParams(q url.Values) dto.QueryParams {
	p := dto.QueryParams{}

	if start := q.Get("start_date"); start != "" {
		p.StartDate = start
	}
	if end := q.Get("end_date"); end != "" {
		p.EndDate = end
	}
	if status := q.Get("status"); status != "" {
		p.Status = status
	}

	// Handle pagination with defaults
	page := 1
	if pageStr := q.Get("page"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil {
			page = val
		}
	}
	p.Page = page

	limit := 10
	if limitStr := q.Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}
	p.Limit = limit

	return p
}
