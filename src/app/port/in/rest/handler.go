package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/carddemo/project/src/app/shared/dto"
)

// Handler aggregates all application HTTP handlers.
type Handler struct {
	TransactionHandler     *TransactionHandler
	BatchSettlementHandler *BatchSettlementHandler
}

// NewHandler constructs a new Handler with dependencies.
func NewHandler(
	txnSvc TransactionService,
	bsSvc BatchSettlementService,
) *Handler {
	return &Handler{
		TransactionHandler:     NewTransactionHandler(txnSvc),
		BatchSettlementHandler: NewBatchSettlementHandler(bsSvc),
	}
}

// RegisterRoutes registers all API routes on the provided router.
func (h *Handler) RegisterRoutes(r *chi.Mux) {
	// Transaction Routes
	r.Post("/transactions", h.TransactionHandler.CreateTransaction)
	r.Get("/transactions/{id}", h.TransactionHandler.GetTransaction)
	r.Get("/accounts/{id}/transactions", h.TransactionHandler.GetAccountTransactions)
	r.Post("/transactions/{id}/void", h.TransactionHandler.VoidTransaction)

	// Batch Settlement Routes
	r.Post("/settlements", h.BatchSettlementHandler.CreateSettlement)
	r.Get("/settlements/{id}", h.BatchSettlementHandler.GetSettlement)
	r.Get("/settlements", h.BatchSettlementHandler.ListSettlements)
}
