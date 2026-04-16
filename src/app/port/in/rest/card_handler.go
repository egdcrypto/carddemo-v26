package rest

import (
	"encoding/json"
	"net/http"

	"github.com/carddemo/project/src/app/card/dto"
	"github.com/go-chi/chi/v5"
)

// CardHandlers handles HTTP requests for Cards
type CardHandlers struct {
	cardRepo    CardRepository
	accountRepo AccountRepository
}

func (h *CardHandlers) IssueCard(w http.ResponseWriter, r *http.Request) {
	// Implementation required
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *CardHandlers) GetCard(w http.ResponseWriter, r *http.Request) {
	// Implementation required
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *CardHandlers) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	// Implementation required
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *CardHandlers) ActivateCard(w http.ResponseWriter, r *http.Request) {
	// Implementation required
	w.WriteHeader(http.StatusNotImplemented)
}
