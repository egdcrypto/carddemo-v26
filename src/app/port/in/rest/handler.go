package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Dependencies are interfaces required by the HTTP handlers
type Dependencies struct {
	CardRepo       CardRepository
	CardPolicyRepo CardPolicyRepository
	AccountRepo    AccountRepository
}

// NewHandler creates the main HTTP router
func NewHandler(deps *Dependencies) http.Handler {
	r := chi.NewRouter()

	// Card Handlers
	h := &CardHandlers{
		cardRepo:    deps.CardRepo,
		accountRepo: deps.AccountRepo,
	}
	r.Post("/accounts/{id}/cards", h.IssueCard)
	r.Get("/cards/{id}", h.GetCard)
	r.Put("/cards/{id}/status", h.UpdateStatus)
	r.Post("/cards/{id}/activate", h.ActivateCard)

	// Policy Handlers
	ph := &PolicyHandlers{
		cardPolicyRepo: deps.CardPolicyRepo,
		accountRepo:    deps.AccountRepo,
	}
	r.Get("/policies/{id}", ph.GetPolicy)
	r.Put("/policies/{id}", ph.UpdatePolicy)
	r.Get("/accounts/{id}/policies", ph.GetAccountPolicies)

	return r
}
