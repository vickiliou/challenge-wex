package httpresponse

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response.
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// RespondWithError returns an error response with the specified status code and error message.
func RespondWithError(w http.ResponseWriter, statusCode int, err error) {
	body := &ErrorResponse{
		StatusCode: statusCode,
		Message:    err.Error(),
	}

	RespondJSON(w, statusCode, body)
}

// RespondJSON returns a JSON response with the specified status code and data.
func RespondJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		LogError("Error encoding response", http.StatusInternalServerError, err)
		return
	}
}
