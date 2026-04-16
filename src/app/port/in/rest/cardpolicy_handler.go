package rest

import (
	"net/http"

	"github.com/carddemo/project/src/app/cardpolicy/dto"
	"github.com/go-chi/chi/v5"
)

// PolicyHandlers handles HTTP requests for CardPolicies
type PolicyHandlers struct {
	cardPolicyRepo CardPolicyRepository
	accountRepo    AccountRepository
}

func (h *PolicyHandlers) GetPolicy(w http.ResponseWriter, r *http.Request) {
	// Implementation required
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *PolicyHandlers) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	// Implementation required
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *PolicyHandlers) GetAccountPolicies(w http.ResponseWriter, r *http.Request) {
	// Implementation required
	w.WriteHeader(http.StatusNotImplemented)
}
