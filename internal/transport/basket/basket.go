package basket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IBasketUsecase interface{
	GetProducts(ctx context.Context)(*models.BasketResponse, error)
	AddProduct(ctx context.Context, productID uuid.UUID)(*models.BasketItem, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID)(error)
	UpdateProductQuantity(ctx context.Context, productID uuid.UUID, quantity int)(*models.BasketItem, error)
	ClearBasket(ctx context.Context)(error)
}

type BasketHandler struct {
	u            IBasketUsecase
	log          *logrus.Logger
}

func NewBasketHandler(u IBasketUsecase, log *logrus.Logger) *BasketHandler {
	return &BasketHandler{
		u:            u,
		log:          log,
	}
}

type addProductRequest struct {
	ProductID uuid.UUID `json:"product_id"`
}

type delProductRequest struct {
	ProductId uuid.UUID `json:"product_id"`
}

type updateQuantityRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
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
func (h *BasketHandler) GetBasket(w http.ResponseWriter, r *http.Request) {
	items, err := h.u.GetProducts(r.Context())
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
//	@Accept			json
//	@Produce		json
//	@Param			request	body		addProductRequest		true	"Данные товара"
//	@Success		201		{object}	models.BasketItem		"Добавленный товар"
//	@Failure		400		{object}	response.ErrorResponse	"Некорректные данные"
//	@Failure		401		{object}	response.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404		{object}	response.ErrorResponse	"Товар не найден"
//	@Failure		500		{object}	response.ErrorResponse	"Ошибка сервера"
//	@Router			/basket/add [post]
func (h *BasketHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var req addProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.u.AddProduct(r.Context(), req.ProductID)
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
//	@Summary	Удалить товар из корзины
//	@Tags		basket
//	@Security	ApiKeyAuth
//	@Param		request	body	delProductRequest	true	"Данные товара"
//	@Success	204		"Товар успешно удалён"
//	@Failure	400		{object}	response.ErrorResponse	"Некорректный формат запроса"
//	@Failure	401		{object}	response.ErrorResponse	"Пользователь не авторизован"
//	@Failure	404		{object}	response.ErrorResponse	"Товар не найден в корзине"
//	@Failure	500		{object}	response.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router		/basket/remove [delete]
func (h *BasketHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var req delProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.u.DeleteProduct(r.Context(), req.ProductId); err != nil {
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
//	@Param			request	body		updateQuantityRequest	true	"Данные для обновления"
//	@Success		200		{object}	models.BasketItem		"Обновленный товар"
//	@Failure		400		{object}	response.ErrorResponse	"Некорректные данные"
//	@Failure		401		{object}	response.ErrorResponse	"Пользователь не авторизован"
//	@Failure		404		{object}	response.ErrorResponse	"Товар не найден"
//	@Failure		500		{object}	response.ErrorResponse	"Ошибка сервера"
//	@Router			/basket/update [put]
func (h *BasketHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	var req updateQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "Некорректный формат данных")
		return
	}

	if req.Quantity <= 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "Количество должно быть положительным")
		return
	}

	item, err := h.u.UpdateProductQuantity(r.Context(), req.ProductID, req.Quantity)
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
func (h *BasketHandler) ClearBasket(w http.ResponseWriter, r *http.Request) {
	if err := h.u.ClearBasket(r.Context()); err != nil {
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