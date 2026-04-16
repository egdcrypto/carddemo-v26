package shared

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// Validate is a singleton validator instance.
	Validate = validator.New()
)

// RespondWithError writes a JSON error response to the writer.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Clean up the error message if it comes from validator.Struct
	// e.g. "Key: 'CreateAccountRequest.UserProfileID' Error:Field validation for 'UserProfileID' failed on the 'required' tag"
	// We want a cleaner message for the API consumer.
	errorResp := map[string]string{"error": simplifyError(message)}
	json.NewEncoder(w).Encode(errorResp)
}

// simplifyError attempts to clean up validator error strings.
func simplifyError(err string) string {
	// Check for validation error format
	if strings.Contains(err, "Error:Field validation") {
		parts := strings.Split(err, "'")
		if len(parts) > 5 {
			field := parts[3]
			tag := parts[len(parts)-1]
			return field + " " + tag
		}
	}
	return err
}
