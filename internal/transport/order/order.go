package order

import (
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type OrderHandler struct {
	u   order.IOrderUsecase
	log *logrus.Logger
}

func NewOrderHandler(
	u order.IOrderUsecase,
	log *logrus.Logger,
) *OrderHandler {
	return &OrderHandler{
		u:   u,
		log: log,
	}
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userIDStr, isExist := r.Context().Value(cookie.UserIDKey).(string)
	if !isExist {
		response.HandleError(w, errs.ErrUserNotFound)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.HandleError(w, errs.ErrInvalidUserID)
		return
	}

	var request models.CreateOrderDTO
	if errStatusCode, err := response.ParseData(r.Body, &request); err != nil {
		response.SendErrorResponse(w, errStatusCode, fmt.Sprintf("Failed to parse request body: %v", err))
		return
	}

	request.UserID = userID
	if err := o.u.CreateOrder(r.Context(), request); err != nil {
		response.HandleError(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, nil)
}
