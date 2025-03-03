package utils

import (
	"encoding/json"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response, err := json.Marshal(ErrorResponse{Message: message})
	if err != nil {
		return
	}
	_, _ = w.Write(response)
}

func SendSuccessResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	if body != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(statusCode)
	if body != nil {
		response, err := json.Marshal(body)
		if err != nil {
			return
		}
		_, _ = w.Write(response)
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
