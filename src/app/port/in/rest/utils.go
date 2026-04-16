package rest

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/carddemo/project/src/app/card/dto"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func validateRequest(req interface{}) error {
	if err := validate.Struct(req); err != nil {
		// Return a generic error or handle field errors specifically
		return err
	}
	return nil
}

// Extension to support validation for structs
func init() {
	// Register custom validation if needed
}
