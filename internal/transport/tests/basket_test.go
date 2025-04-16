package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/basket"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupTestBasket(t *testing.T) (*mocks.MockIBasketUsecase, *basket.BasketService) {
	ctrl := gomock.NewController(t)
	mockUsecase := mocks.NewMockIBasketUsecase(ctrl)
	service := basket.NewBasketService(mockUsecase)
	return mockUsecase, service
}

func TestBasketService_Get(t *testing.T) {
	mockUsecase, service := setupTestBasket(t)

	t.Run("success", func(t *testing.T) {
		expectedItems := []*models.BasketItem{
			{
				ID:            uuid.New(),
				ProductID:     uuid.New(),
				Quantity:      2,
				Price:         10.5,
				PriceDiscount: 9.5,
			},
			{
				ID:        uuid.New(),
				ProductID: uuid.New(),
				Quantity:  1,
				Price:     15.0,
			},
		}

		mockUsecase.EXPECT().
			Get(gomock.Any()).
			Return(expectedItems, nil)

		req := httptest.NewRequest("GET", "/api/v1/basket", nil)
		w := httptest.NewRecorder()

		service.Get(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseData dto.BasketResponse
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)

		// Проверяем правильность преобразования в DTO
		assert.Equal(t, len(expectedItems), responseData.Total)
		assert.Equal(t, 36.0, responseData.TotalPrice)         // 10.5*2 + 15.0*1
		assert.Equal(t, 34.0, responseData.TotalPriceDiscount) // 9.5*2 + 15.0*1
		assert.Equal(t, *expectedItems[0], responseData.Products[0])
		assert.Equal(t, *expectedItems[1], responseData.Products[1])
	})

	t.Run("empty basket", func(t *testing.T) {
		mockUsecase.EXPECT().
			Get(gomock.Any()).
			Return([]*models.BasketItem{}, nil)

		req := httptest.NewRequest("GET", "/api/v1/basket", nil)
		w := httptest.NewRecorder()

		service.Get(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseData dto.BasketResponse
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)
		assert.Equal(t, 0, responseData.Total)
		assert.Equal(t, 0.0, responseData.TotalPrice)
		assert.Equal(t, 0.0, responseData.TotalPriceDiscount)
		assert.Empty(t, responseData.Products)
	})
}

func TestBasketService_Add(t *testing.T) {
	mockUsecase, service := setupTestBasket(t)
	productID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedItem := &models.BasketItem{
			ID:        uuid.New(),
			ProductID: productID,
			Quantity:  1,
			Price:     10.99,
		}

		mockUsecase.EXPECT().
			Add(gomock.Any(), productID).
			Return(expectedItem, nil)

		req := httptest.NewRequest("POST", "/api/v1/basket/"+productID.String(), nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.Add(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseData models.BasketItem
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)
		assert.Equal(t, *expectedItem, responseData)
	})
}

func TestBasketService_UpdateQuantity(t *testing.T) {
	mockUsecase, service := setupTestBasket(t)
	productID := uuid.New()

	t.Run("success", func(t *testing.T) {
		quantity := 3
		expectedItem := &models.BasketItem{
			ID:        uuid.New(),
			ProductID: productID,
			Quantity:  quantity,
			Price:     12.5,
		}

		// Мокируем вызов usecase
		mockUsecase.EXPECT().
			UpdateQuantity(gomock.Any(), productID, quantity).
			Return(expectedItem, nil)

		// Подготавливаем запрос
		requestBody := dto.UpdateQuantityRequest{Quantity: quantity}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("PATCH", "/api/v1/basket/"+productID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Устанавливаем переменные маршрута
		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		// Вызываем обработчик
		service.UpdateQuantity(w, req)

		// Проверяем результат
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Декодируем ответ и проверяем его
		var responseData dto.UpdateQuantityResponse
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)

		// Проверяем, что ответ соответствует ожиданиям
		expectedResponse := dto.ConvertToQuantityResponse(expectedItem)
		assert.Equal(t, expectedResponse, responseData)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("PATCH", "/api/v1/basket/"+productID.String(), strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.UpdateQuantity(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("usecase error", func(t *testing.T) {
		quantity := 2
		mockUsecase.EXPECT().
			UpdateQuantity(gomock.Any(), productID, quantity).
			Return(nil, errs.NewNotFoundError("product not found")) // ✅ теперь 2 аргумента

		requestBody := dto.UpdateQuantityRequest{Quantity: quantity}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("PATCH", "/api/v1/basket/"+productID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.UpdateQuantity(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestConvertToBasketResponse(t *testing.T) {
	t.Run("with discount prices", func(t *testing.T) {
		items := []*models.BasketItem{
			{
				ID:            uuid.New(),
				ProductID:     uuid.New(),
				Quantity:      2,
				Price:         10.0,
				PriceDiscount: 8.0,
			},
			{
				ID:        uuid.New(),
				ProductID: uuid.New(),
				Quantity:  1,
				Price:     15.0,
			},
		}

		result := dto.ConvertToBasketResponse(items)

		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 35.0, result.TotalPrice)         // 10*2 + 15*1
		assert.Equal(t, 31.0, result.TotalPriceDiscount) // 8*2 + 15*1
		assert.Equal(t, *items[0], result.Products[0])
		assert.Equal(t, *items[1], result.Products[1])
	})

	t.Run("empty basket", func(t *testing.T) {
		items := []*models.BasketItem{}
		result := dto.ConvertToBasketResponse(items)

		assert.Equal(t, 0, result.Total)
		assert.Equal(t, 0.0, result.TotalPrice)
		assert.Equal(t, 0.0, result.TotalPriceDiscount)
		assert.Empty(t, result.Products)
	})
}

func TestBasketService_Clear(t *testing.T) {
	mockUsecase, service := setupTestBasket(t)

	t.Run("success", func(t *testing.T) {
		mockUsecase.EXPECT().
			Clear(gomock.Any()).
			Return(nil)

		req := httptest.NewRequest("DELETE", "/basket", nil)
		w := httptest.NewRecorder()

		service.Clear(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockUsecase.EXPECT().
			Clear(gomock.Any()).
			Return(errs.ErrInvalidToken)

		req := httptest.NewRequest("DELETE", "/basket", nil)
		w := httptest.NewRecorder()

		service.Clear(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		mockUsecase.EXPECT().
			Clear(gomock.Any()).
			Return(errors.New("internal error"))

		req := httptest.NewRequest("DELETE", "/basket", nil)
		w := httptest.NewRecorder()

		service.Clear(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestBasketService_Delete(t *testing.T) {
	mockUsecase, service := setupTestBasket(t)
	productID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockUsecase.EXPECT().
			Delete(gomock.Any(), productID).
			Return(nil)

		req := httptest.NewRequest("DELETE", "/api/v1/basket/"+productID.String(), nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("invalid product id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/basket/invalid", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

		service.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("product not found", func(t *testing.T) {
		mockUsecase.EXPECT().
			Delete(gomock.Any(), productID).
			Return(errs.NewNotFoundError("product not found"))

		req := httptest.NewRequest("DELETE", "/api/v1/basket/"+productID.String(), nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockUsecase.EXPECT().
			Delete(gomock.Any(), productID).
			Return(errs.ErrInvalidToken)

		req := httptest.NewRequest("DELETE", "/api/v1/basket/"+productID.String(), nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		mockUsecase.EXPECT().
			Delete(gomock.Any(), productID).
			Return(errors.New("internal error"))

		req := httptest.NewRequest("DELETE", "/api/v1/basket/"+productID.String(), nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
