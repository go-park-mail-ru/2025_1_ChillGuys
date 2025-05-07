package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/address"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест для получения всех адресов пользователя
func TestGetAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок репозитория и юзкейса
	addressUsecase := mocks.NewMockIAddressUsecase(ctrl)
	addressHandler := address.NewAddressHandler(addressUsecase, "geoapify-api-key")

	// Данные для теста
	userID := uuid.New()
	addresses := []dto.GetAddressResDTO{
		{
			ID:            uuid.New(),
			Label:         null.StringFrom("Home"),
			AddressString: null.StringFrom("123 Main St"),
			Coordinate:    null.StringFrom("40.7128, -74.0060"),
		},
	}

	// Мокаем вызовы в юзкейсе
	addressUsecase.EXPECT().GetAddresses(gomock.Any(), userID).Return(addresses, nil).Times(1)

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/addresses", nil)
	req = req.WithContext(context.WithValue(req.Context(), domains.UserIDKey{}, userID.String()))

	// Создаем рекорд
	rec := httptest.NewRecorder()

	// Вызываем хендлер
	addressHandler.GetAddress(rec, req)

	// Проверяем статус и тело
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем ответ
	var responseBody map[string][]dto.GetAddressResDTO
	err := json.NewDecoder(rec.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Len(t, responseBody["addresses"], 1)

	// Проверяем правильность данных
	address := responseBody["addresses"][0]
	assert.Equal(t, addresses[0].ID, address.ID)
	assert.Equal(t, addresses[0].Label, address.Label)
	assert.Equal(t, addresses[0].AddressString, address.AddressString)
	assert.Equal(t, addresses[0].Coordinate, address.Coordinate)
}

func TestGetPickupPoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок репозитория и юзкейса
	addressUsecase := mocks.NewMockIAddressUsecase(ctrl)
	addressHandler := address.NewAddressHandler(addressUsecase, "geoapify-api-key")

	// Данные для теста
	points := []dto.GetPointAddressResDTO{
		{
			ID:            uuid.New(),
			AddressString: null.StringFrom("Pickup Point Address"),
			Coordinate:    null.StringFrom("40.7128, -74.0060"),
		},
	}

	// Мокаем вызовы в юзкейсе
	addressUsecase.EXPECT().GetPickupPoints(gomock.Any()).Return(points, nil).Times(1)

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/addresses/pickup-points", nil)

	// Создаем рекорд
	rec := httptest.NewRecorder()

	// Вызываем хендлер
	addressHandler.GetPickupPoints(rec, req)

	// Проверяем статус и тело
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем ответ
	var responseBody map[string][]dto.GetPointAddressResDTO
	err := json.NewDecoder(rec.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Len(t, responseBody["pickupPoints"], 1)

	// Проверяем правильность данных
	point := responseBody["pickupPoints"][0]
	assert.Equal(t, points[0].ID, point.ID)
	assert.Equal(t, points[0].AddressString, point.AddressString)
	assert.Equal(t, points[0].Coordinate, point.Coordinate)
}

func TestCreateAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIAddressUsecase(ctrl)
	handler := address.NewAddressHandler(mockUsecase, "test-api-key")

	userID := uuid.New()
	addressReq := dto.AddressReqDTO{
		Label:         null.StringFrom("Home"),
		AddressString: null.StringFrom("123 Main St"),
	}

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewReader([]byte("invalid")))
		req = req.WithContext(context.WithValue(req.Context(), domains.UserIDKey{}, userID.String()))

		rec := httptest.NewRecorder()
		handler.CreateAddress(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing user id", func(t *testing.T) {
		body, _ := json.Marshal(addressReq)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewReader(body))

		rec := httptest.NewRecorder()
		handler.CreateAddress(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid user id", func(t *testing.T) {
		body, _ := json.Marshal(addressReq)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), domains.UserIDKey{}, "invalid"))

		rec := httptest.NewRecorder()
		handler.CreateAddress(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func mockResponseBody(data string) *mockReadCloser {
	return &mockReadCloser{
		Reader: bytes.NewReader([]byte(data)),
	}
}

type mockReadCloser struct {
	*bytes.Reader
}

func (m *mockReadCloser) Close() error {
	return nil
}
