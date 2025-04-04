package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ParseData(ioBody io.Reader, request any) (int, error) {
	body, err := io.ReadAll(ioBody)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err)
	}

	return 0, nil
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

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrInvalidCredentials):
		SendErrorResponse(w, http.StatusUnauthorized, "invalid email or password")
	case errors.Is(err, errs.ErrUserNotFound):
		SendErrorResponse(w, http.StatusUnauthorized, "user not found")
	case errors.Is(err, errs.ErrUserAlreadyExists):
		SendErrorResponse(w, http.StatusConflict, "user already exists")
	case errors.Is(err, errs.ErrInvalidUserID):
		SendErrorResponse(w, http.StatusBadRequest, "invalid user id format")
	case errors.Is(err, errs.ErrInvalidToken):
		SendErrorResponse(w, http.StatusUnauthorized, "invalid token")
	case errors.Is(err, errs.ErrProductNotFound):
		SendErrorResponse(w, http.StatusNotFound, "product not found")
	case errors.Is(err, errs.ErrProductNotApproved):
		SendErrorResponse(w, http.StatusForbidden, "product not approved")
	case errors.Is(err, errs.ErrNotEnoughStock):
		SendErrorResponse(w, http.StatusBadRequest, "not enough stock")
	case errors.Is(err, errs.ErrProductDiscountNotFound):
		SendErrorResponse(w, http.StatusNotFound, "product discount not found")
	default:
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}
