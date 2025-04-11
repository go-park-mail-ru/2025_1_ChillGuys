package order

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
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

//	@Summary		Create new order
//	@Description	Создает новый заказ для пользователя
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.CreateOrderDTO	true	"Данные для создания заказа"
//	@Success		200		{}			"Order successfully created"
//	@Failure		400		{object}	response.ErrorResponse	"Некорректный запрос"
//	@Failure		401		{object}	response.ErrorResponse	"Пользователь не найден в контексте"
//	@Failure		404		{object}	response.ErrorResponse	"Ошибка при создании заказа"
//	@Failure		500		{object}	response.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/order [post]

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

	var createOrderReq dto.CreateOrderDTO
	if err := request.ParseData(r, &createOrderReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	createOrderReq.UserID = userID
	if err = o.u.CreateOrder(r.Context(), createOrderReq); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to create order")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

func (o *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userIDStr, isExist := r.Context().Value(domains.UserIDKey).(string)
	if !isExist {
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid user id format")
		return
	}

	orders, err := o.u.GetUserOrders(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to get orders")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string]*[]models.OrderPreview{
		"orders": orders,
	})
}
