package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/app/account/service"
	"github.com/carddemo/project/src/app/shared"
	"github.com/carddemo/project/src/domain/account/command"
	"github.com/go-chi/chi/v5"
)

// AccountHandler handles HTTP requests for Accounts.
type AccountHandler struct {
	service *service.AccountService
}

// NewAccountHandler creates a new REST handler for accounts.
func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// RegisterRoutes mounts the account routes on the router.
func (h *AccountHandler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.CreateAccount)
		r.Get("/{id}", h.GetAccount)
		r.Put("/{id}", h.UpdateAccount)
		r.Delete("/{id}", h.DeleteAccount)

		// Nested routes for UserProfile (handled by UserProfileHandler)
		// We need a reference to UserProfileHandler to mount this properly in a real app,
		// but for the sake of file separation, we'll handle the wiring in the main router setup
		// or accept a sub-handler here. The tests verify specific paths.
	})
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate
	if err := shared.Validate.Struct(req); err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := command.OpenAccountCmd{
		UserProfileID: req.UserProfileID,
		InitialStatus:  req.InitialStatus,
		AccountType:    req.AccountType,
	}
	// Default status if not provided
	if cmd.InitialStatus == "" {
		cmd.InitialStatus = "ACTIVE"
	}

	agg, err := h.service.CreateAccount(cmd)
	if err != nil {
		// Check for specific domain errors if necessary
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	src := dto.MapToAccountResponse(agg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(src)
}

// GetAccount handles GET /accounts/{id}
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	agg, err := h.service.GetAccount(id)
	if err != nil {
		if err == service.ErrAccountNotFound {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.MapToAccountResponse(agg)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateAccount handles PUT /accounts/{id}
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := shared.Validate.Struct(req); err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := command.UpdateAccountStatusCmd{
		NewStatus: req.Status,
		Reason:    req.Reason,
	}

	agg, err := h.service.UpdateAccountStatus(id, cmd)
	if err != nil {
		if err == service.ErrAccountNotFound {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
		// Could be 409 Conflict if state transition invalid, but 500 is safe generic
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.MapToAccountResponse(agg)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteAccount handles DELETE /accounts/{id}
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteAccount(id); err != nil {
		if err == service.ErrAccountNotFound {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
