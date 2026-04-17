package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/carddemo/project/src/app/account/dto"
	"github.com/carddemo/project/src/app/account/service"
	userprofile_dto "github.com/carddemo/project/src/app/userprofile/dto"
	"github.com/carddemo/project/src/app/userprofile/service"
)

// AccountHandler handles HTTP requests.
type AccountHandler struct {
	accountService    *service.AccountService
	userprofileService *service.UserProfileService
	validate          *validator.Validate
}

// NewAccountHandler creates a new handler.
func NewAccountHandler(accService *service.AccountService, profService *service.UserProfileService) *AccountHandler {
	return &AccountHandler{
		accountService:    accService,
		userprofileService: profService,
		validate:          validator.New(),
	}
}

// Routes defines the chi router routes.
func (h *AccountHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateAccount)
	r.Get("/{id}", h.GetAccount)
	r.Put("/{id}", h.UpdateAccount)
	r.Delete("/{id}", h.DeleteAccount)

	r.Get("/{id}/profile", h.GetProfile)
	r.Put("/{id}/profile", h.UpdateProfile)

	return r
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondValidationError(w, err)
		return
	}

	resp, err := h.accountService.CreateAccount(req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusCreated, resp)
}

// GetAccount handles GET /accounts/{id}
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := h.accountService.GetAccount(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "account not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, resp)
}

// UpdateAccount handles PUT /accounts/{id}
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateAccountStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondValidationError(w, err)
		return
	}

	resp, err := h.accountService.UpdateStatus(id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, resp)
}

// DeleteAccount handles DELETE /accounts/{id}
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.accountService.DeleteAccount(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetProfile handles GET /accounts/{id}/profile
func (h *AccountHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	resp, err := h.userprofileService.GetProfileByAccountID(accountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "profile not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, resp)
}

// UpdateProfile handles PUT /accounts/{id}/profile
func (h *AccountHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	var req userprofile_dto.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondValidationError(w, err)
		return
	}

	resp, err := h.userprofileService.UpdateProfileByAccountID(accountID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, resp)
}

// respondJSON writes a JSON response.
func (h *AccountHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// respondValidationError formats validation errors.
func (h *AccountHandler) respondValidationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
}
