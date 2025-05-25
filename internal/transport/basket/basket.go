package basket

import (
	"context"
	"github.com/mailru/easyjson"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=basket.go -destination=../../usecase/mocks/basket_usecase_mock.go -package=mocks IBasketUsecase
type IBasketUsecase interface {
	Get(ctx context.Context) ([]*models.BasketItem, error)
	Add(ctx context.Context, productID uuid.UUID) (*models.BasketItem, error)
	Delete(ctx context.Context, productID uuid.UUID) error
	UpdateQuantity(ctx context.Context, productID uuid.UUID, quantity int) (*models.BasketItem, error)
	Clear(ctx context.Context) error
}

type BasketService struct {
	u IBasketUsecase
}

func NewBasketService(u IBasketUsecase) *BasketService {
	return &BasketService{
		u: u,
	}
}

// GetBasket godoc
//
//	@Summary		Получить содержимое корзины
//	@Description	Возвращает все товары в корзине пользователя
//	@Tags			basket
//	@Produce		json
//	@Success		200	{object}	dto.BasketResponse
//	@Failure		401	{object}	object
//	@Failure		404	{object}	object
//	@Failure		500	{object}	object
//	@Security		TokenAuth
//	@Router			/basket [get]
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
//	@Description	Добавляет товар в корзину пользователя
//	@Tags			basket
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string	true	"ID товара в формате UUID"
//	@Param			X-Csrf-Token	header		string	true	"CSRF-токен для защиты от подделки запросов"
//	@Success		201				{object}	models.BasketItem
//	@Failure		400				{object}	object
//	@Failure		401				{object}	object
//	@Failure		404				{object}	object
//	@Failure		500				{object}	object
//	@Security		TokenAuth
//	@Router			/basket/{id} [post]
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
//	@Param			id				path	string	true	"ID товара в формате UUID"
//	@Param			X-Csrf-Token	header	string	true	"CSRF-токен для защиты от подделки запросов"
//	@Success		204
//	@Failure		400	{object}	object
//	@Failure		401	{object}	object
//	@Failure		404	{object}	object
//	@Failure		500	{object}	object
//	@Security		TokenAuth
//	@Router			/basket/{id} [delete]
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

// UpdateQuantity godoc
//
//	@Summary		Обновить количество товара
//	@Description	Изменяет количество указанного товара в корзине
//	@Tags			basket
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string						true	"ID товара в формате UUID"
//	@Param			request			body		dto.UpdateQuantityRequest	true	"Новое количество"
//	@Param			X-Csrf-Token	header		string						true	"CSRF-токен для защиты от подделки запросов"
//	@Success		200				{object}	dto.UpdateQuantityResponse
//	@Failure		400				{object}	object
//	@Failure		401				{object}	object
//	@Failure		404				{object}	object
//	@Failure		500				{object}	object
//	@Security		TokenAuth
//	@Router			/basket/{id} [patch]
func (h *BasketService) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	const op = "BasketService.UpdateQuantity"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var req dto.UpdateQuantityRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
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
//	@Param			X-Csrf-Token	header	string	true	"CSRF-токен для защиты от подделки запросов"
//	@Success		204
//	@Failure		401	{object}	object
//	@Failure		500	{object}	object
//	@Security		TokenAuth
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
