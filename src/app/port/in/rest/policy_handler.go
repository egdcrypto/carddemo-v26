package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carddemo/project/src/app/card/dto"
	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/card/repository"
	cardpolicycmd "github.com/carddemo/project/src/domain/cardpolicy/command"
	"github.com/carddemo/project/src/domain/cardpolicy/model" // Note: Alias needed as model package name is same as variable
	"github.com/carddemo/project/src/domain/cardpolicy/repository"
	"github.com/go-chi/chi/v5"
)

type PolicyHandler struct {
	policyRepo repository.CardPolicyRepository
	cardRepo   repository.CardRepository
}

func NewPolicyHandler(
	policyRepo repository.CardPolicyRepository,
	cardRepo repository.CardRepository,
) *PolicyHandler {
	return &PolicyHandler{
		policyRepo: policyRepo,
		cardRepo:   cardRepo,
	}
}

func RegisterPolicyRoutes(r chi.Router, pRepo repository.CardPolicyRepository, cRepo repository.CardRepository) {
	h := NewPolicyHandler(pRepo, cRepo)

	r.Get("/policies/{id}", h.GetPolicy)                 // GET /policies/{id}
	r.Put("/policies/{id}", h.UpdatePolicy)              // PUT /policies/{id}
	r.Get("/accounts/{id}/policies", h.GetPoliciesByAccount) // GET /accounts/{id}/policies
}

// GetPolicy handles GET /policies/{id}
func (h *PolicyHandler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	policyID := chi.URLParam(r, "id")

	policy, err := h.policyRepo.Get(policyID)
	if err != nil || policy == nil {
		respondJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: "Policy not found"})
		return
	}

	resp := dto.CardPolicyResponse{
		ID:          policy.ID,
		CardID:      policy.CardID,
		DailyLimit:  policy.DailyLimit,
		WeeklyLimit: policy.WeeklyLimit,
		IsActive:    policy.IsActive,
	}

	respondJSON(w, http.StatusOK, resp)
}

// UpdatePolicy handles PUT /policies/{id}
func (h *PolicyHandler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	policyID := chi.URLParam(r, "id")

	var req dto.UpdateCardPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validateRequest(req); err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	policy, err := h.policyRepo.Get(policyID)
	if err != nil || policy == nil {
		respondJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: "Policy not found"})
		return
	}

	// Update limits using Domain Logic
	cmd := cardpolicycmd.UpdateCardLimitsCmd{
		DailyLimit:  req.DailyLimit,
		WeeklyLimit: req.WeeklyLimit,
	}

	policy.UpdateLimits(cmd)

	if err := h.policyRepo.Save(policy); err != nil {
		respondJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update policy"})
		return
	}

	resp := dto.CardPolicyResponse{
		ID:          policy.ID,
		CardID:      policy.CardID,
		DailyLimit:  policy.DailyLimit,
		WeeklyLimit: policy.WeeklyLimit,
		IsActive:    policy.IsActive,
	}

	respondJSON(w, http.StatusOK, resp)
}

// GetPoliciesByAccount handles GET /accounts/{id}/policies
func (h *PolicyHandler) GetPoliciesByAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	// 1. Get all Cards for the Account
	// Note: The CardRepository interface provided in mocks doesn't have ListByAccount.
	// However, List() is available. We must filter in memory or assume an extended interface.
	// Given the strict rule "DO NOT write new code that exists", and mocks only have List(), 
	// we will fetch all cards and filter. This is inefficient but valid for TDD green phase 
	// with the provided mocks.
	
	allCards, err := h.cardRepo.List()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch cards"})
		return
	}

	var cardIDs []string
	for _, c := range allCards {
		if c.AccountID == accountID {
			cardIDs = append(cardIDs, c.ID)
		}
	}

	// 2. Get all Policies
	allPolicies, err := h.policyRepo.List()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch policies"})
		return
	}

	// 3. Filter policies by Card IDs
	var result []dto.CardPolicyResponse
	for _, p := range allPolicies {
		for _, cid := range cardIDs {
			if p.CardID == cid {
				result = append(result, dto.CardPolicyResponse{
					ID:          p.ID,
					CardID:      p.CardID,
					DailyLimit:  p.DailyLimit,
					WeeklyLimit: p.WeeklyLimit,
					IsActive:    p.IsActive,
				})
				break
			}
		}
	}

	respondJSON(w, http.StatusOK, result)
}
