package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ParseData(ioBody io.Reader, request any) (int, string) {
	body, err := io.ReadAll(ioBody)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	if err := json.Unmarshal(body, &request); err != nil {
		return http.StatusBadRequest, fmt.Sprintf("Failed to parse request body: %v", err.Error())
	}

	return 0, ""
}

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
