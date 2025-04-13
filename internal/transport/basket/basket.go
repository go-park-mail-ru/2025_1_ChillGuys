package basket

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=basket.go -destination=../../usecase/mocks/basket_usecase_mock.go -package=mocks IBasketUsecase
type IBasketUsecase interface{
	Get(ctx context.Context)([]*models.BasketItem, error)
	Add(ctx context.Context, productID uuid.UUID)(*models.BasketItem, error)
	Delete(ctx context.Context, productID uuid.UUID)(error)
	UpdateQuantity(ctx context.Context, productID uuid.UUID, quantity int)(*models.BasketItem, error)
	Clear(ctx context.Context)(error)
}

type BasketService struct {
	u            IBasketUsecase
}

func NewBasketService(u IBasketUsecase) *BasketService {
	return &BasketService{
		u:            u,
	}
}

// GetBasket godoc
//
//	@Summary		Получить содержимое корзины
//	@Description	Возвращает все товары в корзине пользователя
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	dto.BasketResponse	"Содержимое корзины"
//	@Failure		401	{object}	dto.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404	{object}	dto.ErrorResponse	"Корзина не найдена"
//	@Failure		500	{object}	dto.ErrorResponse	"Ошибка сервера"
//	@Router			/api/v1/basket [get]
func (h *BasketService) Get(w http.ResponseWriter, r *http.Request) {
	const op = "BasketService.Get"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	items, err := h.u.Get(r.Context())
	if err != nil {
        logger.WithError(err).Error("get basket items")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	responseBasket := dto.ConvertToBasketResponse(items)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseBasket)
}

// AddToBasket godoc
//
//	@Summary		Добавить товар в корзину
//	@Description	Добавляет товар в корзину пользователя или увеличивает количество, если товар уже есть
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		string				true	"ID товара"
//	@Success		201	{object}	models.BasketItem	"Добавленный товар"
//	@Failure		400	{object}	dto.ErrorResponse	"Некорректный ID"
//	@Failure		401	{object}	dto.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404	{object}	dto.ErrorResponse	"Товар не найден"
//	@Failure		500	{object}	dto.ErrorResponse	"Ошибка сервера"
//	@Router			/api/v1/basket/{id} [post]
func (h *BasketService) Add(w http.ResponseWriter, r *http.Request) {
	const op = "BasketService.Add"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
        logger.WithError(err).WithField("product_id", idStr).Error("parse product ID")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	item, err := h.u.Add(r.Context(), productID)
	if err != nil {
        logger.WithField("product_id", productID).WithError(err).Error("add product to basket")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	response.SendJSONResponse(r.Context(), w, http.StatusCreated, item)
}

// RemoveFromBasket godoc
//
//	@Summary		Удалить товар из корзины
//	@Description	Удаляет товар из корзины пользователя
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Param			id	path	string	true	"ID товара"
//	@Success		204	"Товар успешно удалён"
//	@Failure		400	{object}	dto.ErrorResponse	"Некорректный ID"
//	@Failure		401	{object}	dto.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404	{object}	dto.ErrorResponse	"Товар не найден в корзине"
//	@Failure		500	{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/api/v1/basket/{id} [delete]
func (h *BasketService) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "BasketService.Delete"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
        logger.WithError(err).WithField("product_id", idStr).Error("parse product ID")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }


	if err := h.u.Delete(r.Context(), productID); err != nil {
        logger.WithField("product_id", productID).WithError(err).Error("delete product from basket")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	response.SendJSONResponse(r.Context(), w, http.StatusNoContent, nil)
}

// UpdateBasketItem godoc
//
//	@Summary		Обновить количество товара
//	@Description	Изменяет количество указанного товара в корзине
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string						true	"ID товара"
//	@Param			quantity	body		dto.UpdateQuantityRequest	true	"Новое количество"
//	@Success		200			{object}	models.BasketItem			"Обновленный товар"
//	@Failure		400			{object}	dto.ErrorResponse			"Некорректные данные"
//	@Failure		401			{object}	dto.ErrorResponse			"Пользователь не авторизован"
//	@Failure		404			{object}	dto.ErrorResponse			"Товар не найден"
//	@Failure		500			{object}	dto.ErrorResponse			"Ошибка сервера"
//	@Router			/api/v1/basket/{id} [patch]
func (h *BasketService) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	const op = "BasketService.UpdateQuantity"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var req dto.UpdateQuantityRequest
	if err := request.ParseData(r, &req); err != nil {
        logger.WithError(err).Error("parse request data")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
        logger.WithError(err).WithField("product_id", idStr).Error("parse product ID")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	logger = logger.WithFields(logrus.Fields{
        "product_id": productID,
        "quantity":   req.Quantity,
    })

	item, err := h.u.UpdateQuantity(r.Context(), productID, req.Quantity)
	if err != nil {
        logger.WithError(err).Error("update product quantity")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	resp := dto.ConvertToQuantityResponse(item)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, resp)
}

// ClearBasket godoc
//
//	@Summary		Очистить корзину
//	@Description	Полностью удаляет все товары из корзины пользователя
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Success		204	"Корзина успешно очищена"
//	@Failure		401	{object}	dto.ErrorResponse	"Пользователь не авторизован"
//	@Failure		500	{object}	dto.ErrorResponse	"Ошибка сервера"
//	@Router			/basket [delete]
func (h *BasketService) Clear(w http.ResponseWriter, r *http.Request) {
	const op = "BasketService.Clear"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	if err := h.u.Clear(r.Context()); err != nil {
        logger.WithError(err).Error("clear basket")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	response.SendJSONResponse(r.Context(), w, http.StatusNoContent, nil)
}