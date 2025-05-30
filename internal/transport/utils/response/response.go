package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func SendJSONError(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp, err := json.Marshal(dto.ErrorResponseDTO{Message: message})
	if err != nil {
		logctx.GetLogger(ctx).Error("failed to marshal response: ", err.Error())
		return
	}

	if _, err := w.Write(resp); err != nil {
		logctx.GetLogger(ctx).Error("failed to write response: ", err.Error())
	}
}

func SendJSONResponse(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	if body == nil {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logctx.GetLogger(ctx).Error("failed to marshal response", err.Error())
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(resp); err != nil {
		logctx.GetLogger(ctx).Error("failed to write response", err.Error())
	}
}

func HandleDomainError(ctx context.Context, w http.ResponseWriter, err error, description string) {
	log := logctx.GetLogger(ctx)

	switch {
	case errors.Is(err, errs.ErrInvalidCredentials):
		SendJSONError(ctx, w, http.StatusUnauthorized, fmt.Sprintf("%s: %v", description, err))
		log.Debug("invalid credentials error: ", description, err.Error())

	case errors.Is(err, errs.ErrNotFound):
		SendJSONError(ctx, w, http.StatusUnauthorized, fmt.Sprintf("%s: %v", description, err))
		log.Debug("user not found: ", description, err.Error())

	case errors.Is(err, errs.ErrAlreadyExists):
		SendJSONError(ctx, w, http.StatusConflict, fmt.Sprintf("%s: %v", description, err))
		log.Debug("user already exists: ", description, err.Error())

	case errors.Is(err, errs.ErrInvalidID):
		SendJSONError(ctx, w, http.StatusBadRequest, fmt.Sprintf("%s: %v", description, err))
		log.Debug("invalid user id format: ", description, err.Error())

	case errors.Is(err, errs.ErrInvalidToken):
		SendJSONError(ctx, w, http.StatusUnauthorized, fmt.Sprintf("%s: %v", description, err))
		log.Debug("invalid token: ", description, err.Error())

	case errors.Is(err, errs.ErrAlreadyExists):
		SendJSONError(ctx, w, http.StatusConflict, fmt.Sprintf("%s: %v", description, err))
		log.Debug("resource already exists: ", description, errs.NewAlreadyExistsError(description))

	case errors.Is(err, errs.ErrNotFound):
		SendJSONError(ctx, w, http.StatusNotFound, fmt.Sprintf("%s: %v", description, err))
		log.Debug("not found: ", description, errs.NewNotFoundError(description))

	case errors.Is(err, errs.ErrInvalidID):
		SendJSONError(ctx, w, http.StatusBadRequest, fmt.Sprintf("%s: %v", description, err))
		log.Debug("invalid id format: ", description, err.Error())

	case errors.Is(err, errs.ErrBusinessLogic):
		SendJSONError(ctx, w, http.StatusUnprocessableEntity, fmt.Sprintf("%s: %v", description, err))
		log.Debug("business logic error: ", description, errs.NewBusinessLogicError(description))

	case errors.Is(err, errs.ErrNotFound):
		SendJSONError(ctx, w, http.StatusNotFound, fmt.Sprintf("%s: %v", description, err))
		log.Debug("product not found: ", description, err.Error())

	case errors.Is(err, errs.ErrProductNotApproved):
		SendJSONError(ctx, w, http.StatusForbidden, fmt.Sprintf("%s: %v", description, err))
		log.Debug("product not approved: ", description, err.Error())

	case errors.Is(err, errs.ErrNotEnoughStock):
		SendJSONError(ctx, w, http.StatusBadRequest, fmt.Sprintf("%s: %v", description, err))
		log.Debug("not enough stock: ", description, err.Error())
	
	case errors.Is(err, errs.ErrInvalidProductPrice):
		SendJSONError(ctx, w, http.StatusBadRequest, fmt.Sprintf("%s: %v", description, err))
		log.Debug("invalid format: ", description, err.Error())

	case errors.Is(err, errs.ErrInvalidProductQuantity):
		SendJSONError(ctx, w, http.StatusBadRequest, fmt.Sprintf("%s: %v", description, err))
		log.Debug("invalid format: ", description, err.Error())

	case errors.Is(err, errs.ErrEmptyProductName):
		SendJSONError(ctx, w, http.StatusBadRequest, fmt.Sprintf("%s: %v", description, err))
		log.Debug("invalid format: ", description, err.Error())

	default:
		SendJSONError(ctx, w, http.StatusInternalServerError, err.Error())
		log.Error("unexpected error: ", description, err.Error())
	}
}

func HandleGRPCError(ctx context.Context, w http.ResponseWriter, err error, op string) {
	logger := logctx.GetLogger(ctx)
	st, ok := status.FromError(err)
	if !ok {
		logger.WithError(err).Error(op + ": unexpected error type")
		SendJSONError(ctx, w, http.StatusInternalServerError, "internal server error")
		return
	}

	switch st.Code() {
	case codes.Unauthenticated:
		SendJSONError(ctx, w, http.StatusUnauthorized, st.Message())
	case codes.AlreadyExists:
		SendJSONError(ctx, w, http.StatusConflict, st.Message())
	case codes.NotFound:
		SendJSONError(ctx, w, http.StatusNotFound, st.Message())
	case codes.InvalidArgument:
		SendJSONError(ctx, w, http.StatusBadRequest, st.Message())
	default:
		logger.WithError(err).Error(op + ": unexpected gRPC status code")
		SendJSONError(ctx, w, http.StatusInternalServerError, "internal server error")
	}
}
