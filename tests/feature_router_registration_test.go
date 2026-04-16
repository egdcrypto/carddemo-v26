package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/app/port/in/rest"
	cardpolicyrepo "github.com/carddemo/project/src/domain/cardpolicy/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TestRouter_Registration verifies that route paths are correctly mapped.
func TestRouter_Registration(t *testing.T) {
	mockCardRepo := mocks.NewMockCardRepository()
	mockAccountRepo := mocks.NewMockAccountRepository()
	mockPolicyRepo := mocks.NewMockCardPolicyRepository()

	r := chi.NewRouter()
	rest.RegisterCardRoutes(r, mockCardRepo, mockAccountRepo, mockPolicyRepo)
	rest.RegisterPolicyRoutes(r, mockPolicyRepo, mockCardRepo)

	ts := httptest.NewServer(r)
	defer ts.Close()

	tests := []struct {
		method string
		path   string
		status int
	}{
		{"POST", "/accounts/acc_1/cards", http.StatusCreated}, // Expect Created or Error depending on logic, but route must exist
		{"GET", "/cards/card_1", http.StatusNotFound},          // Returns 404 as mock is empty
		{"PUT", "/cards/card_1/status", http.StatusBadRequest}, // Returns 400 because body is empty/malformed
		{"POST", "/cards/card_1/activate", http.StatusBadRequest},
		{"GET", "/policies/pol_1", http.StatusNotFound},
		{"PUT", "/policies/pol_1", http.StatusBadRequest},
		{"GET", "/accounts/acc_1/policies", http.StatusOK}, // Empty array 200
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, ts.URL+tc.path, nil)
			client := &http.Client{}
			res, err := client.Do(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.status, res.StatusCode, "Route might not be registered correctly")
		})
	}
}
