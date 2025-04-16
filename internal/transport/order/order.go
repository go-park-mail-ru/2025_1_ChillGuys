package order

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OrderService struct {
	u   order.IOrderUsecase
}

func NewOrderService(
	u order.IOrderUsecase,
) *OrderService {
	return &OrderService{
		u:   u,
	}
}

// CreateOrder godoc
//
//	@Summary		Создать новый заказ
//	@Description	Создает новый заказ для текущего пользователя
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			orderData		body	dto.CreateOrderDTO	true	"Данные для создания заказа"
//	@Param			X-Csrf-Token	header	string				true	"CSRF-токен для защиты от подделки запросов"
//	@Success		200				"Заказ успешно создан"
//	@Failure		400				{object}	object	"Некорректные данные"
//	@Failure		401				{object}	object	"Пользователь не авторизован"
//	@Failure		404				{object}	object	"Ошибка при создании заказа"
//	@Failure		500				{object}	object	"Внутренняя ошибка сервера"
//	@Security		TokenAuth
//	@Router			/orders [post]
func (o *OrderService) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const op = "OrderService.CreateOrder"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)
	
	userIDStr, isExist := r.Context().Value(domains.UserIDKey{}).(string)
	if !isExist {
		logger.Error("user not found in context")
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "user not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Error("invalid user id format")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid user id format")
		return
	}

	var createOrderReq dto.CreateOrderDTO
	if err := request.ParseData(r, &createOrderReq); err != nil {
		logger.WithError(err).Error("parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	createOrderReq.UserID = userID
	logger = logger.WithFields(logrus.Fields{
		"user_id": userID,
		"order":   createOrderReq,
	})
	if err = o.u.CreateOrder(r.Context(), createOrderReq); err != nil {
		logger.WithError(err).Error("create order")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// GetOrders godoc
//
//	@Summary		Получить список заказов
//	@Description	Возвращает список всех заказов текущего пользователя
//	@Tags			order
//	@Produce		json
//	@Success		200	{object}	map[string][]dto.OrderPreviewDTO	"Список заказов"
//	@Failure		400	{object}	object								"Некорректный ID пользователя"
//	@Failure		401	{object}	object								"Пользователь не авторизован"
//	@Failure		500	{object}	object								"Внутренняя ошибка сервера"
//	@Security		TokenAuth
//	@Router			/orders [get]
func (o *OrderService) GetOrders(w http.ResponseWriter, r *http.Request) {
	const op = "OrderService.GetOrders"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	userIDStr, isExist := r.Context().Value(domains.UserIDKey{}).(string)
	if !isExist {
		logger.Error("user id not found in context")
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Error("invalid user id format")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid user id format")
		return
	}

	logger = logger.WithField("user_id", userID)
	orders, err := o.u.GetUserOrders(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Error("get user orders")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string]*[]dto.OrderPreviewDTO{
		"orders": orders,
	})
}
