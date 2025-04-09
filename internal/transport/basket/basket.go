package basket

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
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
func (h *BasketService) Get(w http.ResponseWriter, r *http.Request) {
	items, err := h.u.Get(r.Context())
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "get basket")
		return
	}

	responseBasket := dto.ConvertToBasketResponse(items)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseBasket)
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
func (h *BasketService) Add(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		response.HandleDomainError(r.Context(), w, err, "add product in basket")
		return
	}

	item, err := h.u.Add(r.Context(), productID)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "add product in basket")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusCreated, item)
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
func (h *BasketService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		response.HandleDomainError(r.Context(), w, err, "del product of basket")
		return
	}

	if err := h.u.Delete(r.Context(), productID); err != nil {
		response.HandleDomainError(r.Context(), w, err, "del product of basket")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusNoContent, nil)
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
	if err := request.ParseData(r, &req); err != nil {
		response.HandleDomainError(r.Context(), w, err, "update quantity")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		response.HandleDomainError(r.Context(), w, err, "update quantity")
		return
	}

	item, err := h.u.UpdateQuantity(r.Context(), productID, req.Quantity)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "update quantity")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, item)
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
func (h *BasketService) Clear(w http.ResponseWriter, r *http.Request) {
	if err := h.u.Clear(r.Context()); err != nil {
		response.HandleDomainError(r.Context(), w, err, "clear basket")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusNoContent, nil)
}