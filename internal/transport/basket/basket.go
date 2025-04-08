package basket

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=basket.go -destination=../../usecase/mocks/basket_usecase_mock.go -package=mocks IBasketUsecase
type IBasketUsecase interface{
	Get(ctx context.Context)(*dto.BasketResponse, error)
	AddProduct(ctx context.Context, productID uuid.UUID)(*models.BasketItem, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID)(error)
	UpdateProductQuantity(ctx context.Context, productID uuid.UUID, quantity int)(*models.BasketItem, error)
	Clear(ctx context.Context)(error)
}

type BasketService struct {
	u            IBasketUsecase
	log          *logrus.Logger
}

func NewBasketService(u IBasketUsecase, log *logrus.Logger) *BasketService {
	return &BasketService{
		u:            u,
		log:          log,
	}
}

// GetBasket godoc
//
//	@Summary		Получить содержимое корзины
//	@Description	Возвращает все товары в корзине пользователя
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	models.BasketResponse	"Содержимое корзины"
//	@Failure		401	{object}	response.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404	{object}	response.ErrorResponse	"Корзина не найдена"
//	@Failure		500	{object}	response.ErrorResponse	"Ошибка сервера"
//	@Router			/basket [get]
func (h *BasketService) GetBasket(w http.ResponseWriter, r *http.Request) {
	items, err := h.u.Get(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			response.SendErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		case errors.Is(err, errs.ErrNotFound):
			response.SendErrorResponse(w, http.StatusNotFound, "Basket not found")
		default:
			h.log.Errorf("Get basket error: %v", err)
			response.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	response.SendSuccessResponse(w, http.StatusOK, items)
}

// AddProduct godoc
//
//	@Summary		Добавить товар в корзину
//	@Description	Добавляет товар в корзину пользователя или увеличивает количество, если товар уже есть
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		string					true	"ID товара"
//	@Success		201	{object}	models.BasketItem		"Добавленный товар"
//	@Failure		400	{object}	response.ErrorResponse	"Некорректный ID"
//	@Failure		401	{object}	response.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404	{object}	response.ErrorResponse	"Товар не найден"
//	@Failure		500	{object}	response.ErrorResponse	"Ошибка сервера"
//	@Router			/basket/add/{id} [post]
func (h *BasketService) AddProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		response.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	item, err := h.u.AddProduct(r.Context(), productID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			response.SendErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		case errors.Is(err, errs.ErrInvalidProductID), errors.Is(err, errs.ErrNotFound):
			response.SendErrorResponse(w, http.StatusInternalServerError, "Invalid product ID")
		case errors.Is(err, errs.ErrNotFound), errors.Is(err, errs.ErrNotFound):
			response.SendErrorResponse(w, http.StatusNotFound, "Invalid product ID")
		default:
			response.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, item)
}

// DeleteProduct godoc
//
//	@Summary		Удалить товар из корзины
//	@Description	Удаляет товар из корзины пользователя
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Param			id	path		string					true	"ID товара"
//	@Success		204	"Товар успешно удалён"
//	@Failure		400	{object}	response.ErrorResponse	"Некорректный ID"
//	@Failure		401	{object}	response.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404	{object}	response.ErrorResponse	"Товар не найден в корзине"
//	@Failure		500	{object}	response.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/basket/remove/{id} [delete]
func (h *BasketService) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		response.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.u.DeleteProduct(r.Context(), productID); err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			response.SendErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		case errors.Is(err, errs.ErrNotFound):
			response.SendErrorResponse(w, http.StatusNotFound, "Product not found in basket")
		default:
			h.log.Errorf("Delete product error: %v", err)
			response.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	response.SendSuccessResponse(w, http.StatusNoContent, nil)
}

// UpdateQuantity godoc
//
//	@Summary		Обновить количество товара
//	@Description	Изменяет количество указанного товара в корзине
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"ID товара"
//	@Param			quantity	query		int						true	"Новое количество"
//	@Success		200		{object}	models.BasketItem		"Обновленный товар"
//	@Failure		400		{object}	response.ErrorResponse	"Некорректные данные"
//	@Failure		401		{object}	response.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404		{object}	response.ErrorResponse	"Товар не найден"
//	@Failure		500		{object}	response.ErrorResponse	"Ошибка сервера"
//	@Router			/basket/update/{id} [patch]
func (h *BasketService) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateQuantityRequest

	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		response.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	item, err := h.u.UpdateProductQuantity(r.Context(), productID, req.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			response.SendErrorResponse(w, http.StatusUnauthorized, "Требуется авторизация")
		case errors.Is(err, errs.ErrNotFound):
			response.SendErrorResponse(w, http.StatusNotFound, "Товар не найден в корзине")
		case errors.Is(err, errs.ErrInvalidProductID):
			response.SendErrorResponse(w, http.StatusBadRequest, "Некорректный ID товара")
		default:
			h.log.Errorf("Ошибка обновления количества: %v", err)
			response.SendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, item)
}

// ClearBasket godoc
//
//	@Summary		Очистить корзину
//	@Description	Полностью удаляет все товары из корзины пользователя
//	@Tags			basket
//	@Security		ApiKeyAuth
//	@Success		204	"Корзина успешно очищена"
//	@Failure		401	{object}	response.ErrorResponse	"Пользователь не авторизован"
//	@Failure		500	{object}	response.ErrorResponse	"Ошибка сервера"
//	@Router			/basket/clear [delete]
func (h *BasketService) ClearBasket(w http.ResponseWriter, r *http.Request) {
	if err := h.u.Clear(r.Context()); err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			response.SendErrorResponse(w, http.StatusUnauthorized, "Требуется авторизация")
		default:
			h.log.Errorf("Ошибка очистки корзины: %v", err)
			response.SendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	response.SendSuccessResponse(w, http.StatusNoContent, nil)
}