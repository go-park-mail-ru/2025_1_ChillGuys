package order

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
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
	userIDStr, isExist := r.Context().Value(domains.UserIDKey).(string)
	if !isExist {
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "user not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid user id format")
		return
	}

	var CreateOrderReq models.CreateOrderDTO
	if err := request.ParseData(r, &CreateOrderReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	CreateOrderReq.UserID = userID
	if err := o.u.CreateOrder(r.Context(), CreateOrderReq); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to create order")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}
