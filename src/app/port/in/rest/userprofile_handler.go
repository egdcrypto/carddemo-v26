package rest

import (
	"encoding/json"
	"net/http"

	"github.com/carddemo/project/src/app/shared"
	userprofiledto "github.com/carddemo/project/src/app/userprofile/dto"
	"github.com/carddemo/project/src/app/userprofile/service"
	"github.com/carddemo/project/src/domain/userprofile/command"
	"github.com/go-chi/chi/v5"
)

// UserProfileHandler handles HTTP requests for UserProfiles.
type UserProfileHandler struct {
	service *service.UserProfileService
}

// NewUserProfileHandler creates a new REST handler for user profiles.
func NewUserProfileHandler(service *service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{service: service}
}

// RegisterRoutes mounts the user profile routes. Note: In chi, these are nested under accounts.
func (h *UserProfileHandler) RegisterRoutes(r chi.Router) {
	r.Get("/accounts/{accountID}/profile", h.GetProfile)
	r.Put("/accounts/{accountID}/profile", h.UpdateProfile)
}

// GetProfile handles GET /accounts/{accountID}/profile
func (h *UserProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")

	profile, err := h.service.GetProfileByAccountID(accountID)
	if err != nil {
		if err == service.ErrProfileNotFound {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := userprofiledto.MapToUserProfileResponse(profile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateProfile handles PUT /accounts/{accountID}/profile
func (h *UserProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")

	var req userprofiledto.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := shared.Validate.Struct(req); err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// First, we need to find the profile by AccountID to get its actual ID
	profile, err := h.service.GetProfileByAccountID(accountID)
	if err != nil {
		if err == service.ErrProfileNotFound {
			http.Error(w, "Profile not found for this account", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cmd := command.UpdateProfileDetailsCmd{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	updatedProfile, err := h.service.UpdateProfileDetails(profile.ID, cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := userprofiledto.MapToUserProfileResponse(updatedProfile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
