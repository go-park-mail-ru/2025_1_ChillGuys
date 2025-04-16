package order

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type OrderService struct {
	u   order.IOrderUsecase
	log *logrus.Logger
}

func NewOrderService(
	u order.IOrderUsecase,
	log *logrus.Logger,
) *OrderService {
	return &OrderService{
		u:   u,
		log: log,
	}
}

// CreateOrder godoc
//
//	@Summary		Создать новый заказ
//	@Description	Создает новый заказ для текущего пользователя
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			orderData	body	dto.CreateOrderDTO	true	"Данные для создания заказа"
//	@Success		200			"Заказ успешно создан"
//	@Failure		400			{object}	object	"Некорректные данные"
//	@Failure		401			{object}	object	"Пользователь не авторизован"
//	@Failure		404			{object}	object	"Ошибка при создании заказа"
//	@Failure		500			{object}	object	"Внутренняя ошибка сервера"
//	@Security		TokenAuth
//	@Router			/orders [post]
func (o *OrderService) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userIDStr, isExist := r.Context().Value(domains.UserIDKey{}).(string)
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
	userIDStr, isExist := r.Context().Value(domains.UserIDKey{}).(string)
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

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string]*[]dto.OrderPreviewDTO{
		"orders": orders,
	})
}
