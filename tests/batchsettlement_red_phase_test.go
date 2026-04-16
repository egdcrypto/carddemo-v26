package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/carddemo/project/src/app/batchsettlement/dto"
	"github.com/stretchr/testify/assert"
)

func (s *RESTHandlerTestSuite) Test_BatchSettlementEndpoints_Post_Settlement_Success() {
	// TDD RED PHASE:
	// We want 201 Created.

	payload := dto.CreateSettlementRequest{
		MerchantID: "merch_123",
		Currency:   "USD",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusCreated, resp.Code, "Expected 201 Created on settlement creation")
	
	var response dto.SettlementResponse
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(s.T(), err, "Expected valid JSON response")
	assert.NotEmpty(s.T(), response.ID, "Expected non-empty ID")
}

func (s *RESTHandlerTestSuite) Test_BatchSettlementEndpoints_Post_Settlement_ValidationError() {
	// TDD RED PHASE:
	// We want 400 Bad Request.

	payload := dto.CreateSettlementRequest{
		Currency: "INVALID", // Invalid currency length
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code, "Expected 400 Bad Request for invalid currency")
}

func (s *RESTHandlerTestSuite) Test_BatchSettlementEndpoints_Get_Settlement_Success() {
	// TDD RED PHASE:
	// We want 200 OK.

	req := httptest.NewRequest("GET", "/settlements/batch_123", nil)
	resp := httptest.NewRecorder()

	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code, "Expected 200 OK when fetching settlement")
}

func (s *RESTHandlerTestSuite) Test_BatchSettlementEndpoints_Get_Settlements_List_Success() {
	// TDD RED PHASE:
	// We want 200 OK.

	req := httptest.NewRequest("GET", "/settlements", nil)
	resp := httptest.NewRecorder()

	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code, "Expected 200 OK when listing settlements")
}
