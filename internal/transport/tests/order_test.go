package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func addUserIDToContext(r *http.Request, userID uuid.UUID) *http.Request {
	ctx := context.WithValue(r.Context(), domains.UserIDKey{}, userID.String())
	return r.WithContext(ctx)
}

func TestCreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIOrderUsecase(ctrl)
	handler := order.NewOrderService(mockUsecase)

	userID := uuid.New()
	addressID := uuid.New()
	productID := uuid.New()

	reqBody := dto.CreateOrderDTO{
		Items: []dto.CreateOrderItemDTO{
			{
				ProductID: productID,
				Price:     100,
				Quantity:  2,
			},
		},
		AddressID: addressID,
	}
	body, _ := json.Marshal(reqBody)

	r := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = addUserIDToContext(r, userID)

	reqBody.UserID = userID
	mockUsecase.EXPECT().
		CreateOrder(gomock.Any(), reqBody).
		Return(nil)

	w := httptest.NewRecorder()
	handler.CreateOrder(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateOrder_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := order.NewOrderService(nil)

	r := httptest.NewRequest(http.MethodPost, "/orders", nil)
	ctx := context.WithValue(r.Context(), domains.UserIDKey{}, "not-a-uuid")
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateOrder(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrder_NoUserInContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := order.NewOrderService(nil)

	r := httptest.NewRequest(http.MethodPost, "/orders", nil)
	w := httptest.NewRecorder()

	handler.CreateOrder(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetOrders_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIOrderUsecase(ctrl)
	handler := order.NewOrderService(mockUsecase)

	userID := uuid.New()
	expected := &[]dto.OrderPreviewDTO{}

	mockUsecase.EXPECT().
		GetUserOrders(gomock.Any(), userID).
		Return(expected, nil)

	r := httptest.NewRequest(http.MethodGet, "/orders", nil)
	r = addUserIDToContext(r, userID)

	w := httptest.NewRecorder()
	handler.GetOrders(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]*[]dto.OrderPreviewDTO
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp["orders"])
}

func TestGetOrders_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIOrderUsecase(ctrl)
	handler := order.NewOrderService(mockUsecase)

	userID := uuid.New()

	mockUsecase.EXPECT().
		GetUserOrders(gomock.Any(), userID).
		Return(nil, errors.New("db error"))

	r := httptest.NewRequest(http.MethodGet, "/orders", nil)
	r = addUserIDToContext(r, userID)

	w := httptest.NewRecorder()
	handler.GetOrders(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
