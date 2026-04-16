package account

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/go-chi/chi/v5"
)

// AccountHandler handles HTTP requests for the Account aggregate.
type AccountHandler struct {
	repo repository.AccountRepository
}

// NewAccountHandler creates a new handler.
func NewAccountHandler(repo repository.AccountRepository) *AccountHandler {
	return &AccountHandler{repo: repo}
}

// RegisterRoutes mounts the handler's routes onto the router.
func (h *AccountHandler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.CreateAccount)
		r.Get("/{id}", h.GetAccount)
		r.Put("/{id}", h.UpdateAccount)
		r.Delete("/{id}", h.DeleteAccount)
		
		// Sub-routes for profile are technically part of this aggregate context
		// but handled separately if needed. The prompt focuses on /accounts.
		r.Get("/{id}/profile", h.GetUserProfile) // Stub for AC
		r.Put("/{id}/profile", h.UpdateUserProfile) // Stub for AC
	})
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// --- Validation ---
	if req.UserProfileID == "" {
		http.Error(w, "user_profile_id is required", http.StatusBadRequest)
		return
	}
	if req.Status != "Active" && req.Status != "Pending" && req.Status != "Suspended" {
		http.Error(w, "status must be Active, Pending, or Suspended", http.StatusBadRequest)
		return
	}
	if req.AccountType != "Checking" && req.AccountType != "Savings" && req.AccountType != "Credit" {
		http.Error(w, "account_type must be Checking, Savings, or Credit", http.StatusBadRequest)
		return
	}

	// --- Domain Logic ---
	// In a real scenario, we might load UserProfile first.
	// For now, we proceed.

	// Prepare command
	cmd := command.OpenAccountCmd{
		UserProfileID: req.UserProfileID,
		InitialStatus: req.Status,
		AccountType:   req.AccountType,
	}

	// Create Aggregate (Use Case Logic inside handler for simplicity of this example, or Service)
	// Ideally: h.service.OpenAccount(cmd)
	// Here: Direct aggregate invocation
	newAggregate := model.NewAccount(req.UserProfileID, req.Status, req.AccountType)

	// Execute command
	err := newAggregate.Execute(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Persist
	if err := h.repo.Save(newAggregate); err != nil {
		http.Error(w, "failed to save account", http.StatusInternalServerError)
		return
	}

	// Response
	resp := dto.AccountResponse{
		ID:            newAggregate.ID,
		UserProfileID: newAggregate.UserProfileID,
		Status:        newAggregate.Status,
		AccountType:   newAggregate.AccountType,
		Version:       newAggregate.Version,
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

	agg, err := h.repo.Get(id)
	if err != nil || agg == nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	resp := dto.AccountResponse{
		ID:            agg.ID,
		UserProfileID: agg.UserProfileID,
		Status:        agg.Status,
		AccountType:   agg.AccountType,
		Version:       agg.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateAccount handles PUT /accounts/{id}
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateAccountStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.NewStatus == "" {
		http.Error(w, "new_status is required", http.StatusBadRequest)
		return
	}

	// Load Aggregate
	agg, err := h.repo.Get(id)
	if err != nil || agg == nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	cmd := command.UpdateAccountStatusCmd{
		NewStatus: req.NewStatus,
		Reason:    req.Reason,
	}

	if err := agg.Execute(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.repo.Save(agg); err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}

	resp := dto.AccountResponse{
		ID:            agg.ID,
		UserProfileID: agg.UserProfileID,
		Status:        agg.Status,
		AccountType:   agg.AccountType,
		Version:       agg.Version,
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

	// Check existence
	_, err := h.repo.Get(id)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserProfile handles GET /accounts/{id}/profile
// Placeholder to satisfy routing AC for now.
func (h *AccountHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// In a real app, we would query the UserProfile projection/aggregate via UserProfileRepo
	// For now, we return a 404 to indicate the link isn't established yet or logic is missing
	http.Error(w, "UserProfile not implemented yet", http.StatusNotImplemented)
}

// UpdateUserProfile handles PUT /accounts/{id}/profile
func (h *AccountHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "UserProfile not implemented yet", http.StatusNotImplemented)
}
