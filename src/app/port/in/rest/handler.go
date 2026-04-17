package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/app/account/service"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/go-chi/chi/v5"
)

// AccountHandler bundles services and dependencies for HTTP handling.
type AccountHandler struct {
	accountService   *service.AccountApplicationService
	profileService   *service.UserProfileApplicationService
}

// NewAccountHandler creates a new rest handler.
func NewAccountHandler(
	accSvc *service.AccountApplicationService,
	profSvc *service.UserProfileApplicationService,
) *AccountHandler {
	return &AccountHandler{
		accountService: accSvc,
		profileService: profSvc,
	}
}

// RegisterRoutes registers the account routes.
func (h *AccountHandler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.CreateAccount)
		r.Get("/{id}", h.GetAccount)
		r.Delete("/{id}", h.DeleteAccount)
		
		// Profile sub-routes
		r.Put("/{id}/profile", h.UpdateProfile) // Matches PUT /accounts/{id}/profile
	})
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Basic Validation
	if strings.TrimSpace(req.UserProfileID) == "" ||
		strings.TrimSpace(req.AccountType) == "" ||
		strings.TrimSpace(req.Status) == "" {
		http.Error(w, `{"error": "missing required fields"}`, http.StatusBadRequest)
		return
	}

	acc, err := h.accountService.CreateAccount(req.UserProfileID, req.AccountType, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.AccountResponse{
		ID:            acc.ID,
		UserProfileID: acc.UserProfileID,
		Status:        acc.Status,
		AccountType:   acc.AccountType,
		CreatedAt:     acc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     acc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Version:       acc.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetAccount handles GET /accounts/{id}
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	acc, err := h.accountService.GetAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if acc == nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	resp := dto.AccountResponse{
		ID:            acc.ID,
		UserProfileID: acc.UserProfileID,
		Status:        acc.Status,
		AccountType:   acc.AccountType,
		CreatedAt:     acc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     acc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Version:       acc.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteAccount handles DELETE /accounts/{id}
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	err := h.accountService.DeleteAccount(id)
	if err != nil {
		// In mocks, delete on non-existent doesn't error, but if it did:
		if err == model.ErrAccountNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateProfile handles PUT /accounts/{id}/profile
func (h *AccountHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		http.Error(w, "missing account id", http.StatusBadRequest)
		return
	}

	var req dto.LinkUserToAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.FirstName) == "" || strings.TrimSpace(req.LastName) == "" {
		http.Error(w, `{"error": "first_name and last_name required"}`, http.StatusBadRequest)
		return
	}

	prof, err := h.profileService.LinkOrUpdateProfile(accountID, req.FirstName, req.LastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.UserProfileResponse{
		ID:        prof.ID,
		AccountID: prof.AccountID,
		FirstName: prof.FirstName,
		LastName:  prof.LastName,
		Email:     prof.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
