package rest

import (
	"github.com/carddemo/project/src/app/transaction/dto"
	"github.com/carddemo/project/src/domain/transaction/repository"
	"github.com/go-playground/validator/v10"
)

// TransactionHandler handles HTTP requests for Transactions.
type TransactionHandler struct {
	repo     repository.TransactionRepository
	validate *validator.Validate
}

// NewTransactionHandler initializes the handler. FIXED: Corrected signature and body.
func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		repo:     repo,
		validate: validator.New(),
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
	// Implementation pending
}

// GetTransaction handles GET /transactions/{id}
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	// Implementation pending
}

// VoidTransaction handles POST /transactions/{id}/void
func (h *TransactionHandler) VoidTransaction(w http.ResponseWriter, r *http.Request) {
	// Implementation pending
}

// ListAccountTransactions handles GET /accounts/{id}/transactions
func (h *TransactionHandler) ListAccountTransactions(w http.ResponseWriter, r *http.Request) {
	// Implementation pending
}
