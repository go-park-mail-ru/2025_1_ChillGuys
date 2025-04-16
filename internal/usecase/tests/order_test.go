package tests

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIOrderRepository(ctrl)
	logger := logrus.New()

	orderUC := order.NewOrderUsecase(mockRepo, logger)

	testUserID := uuid.New()
	testAddressID := uuid.New()
	testProductID := uuid.New()

	tests := []struct {
		name          string
		input         dto.CreateOrderDTO
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Successful order creation",
			input: dto.CreateOrderDTO{
				UserID:    testUserID,
				AddressID: testAddressID,
				Items: []dto.CreateOrderItemDTO{
					{
						ProductID: testProductID,
						Quantity:  2,
					},
				},
			},
			mockSetup: func() {
				// Mock product price check
				mockRepo.EXPECT().ProductPrice(gomock.Any(), testProductID).
					Return(&models.Product{
						ID:       testProductID,
						Price:    100.0,
						Quantity: 5,
						Status:   models.ProductApproved,
					}, nil)

				// Mock discount check (no discount)
				mockRepo.EXPECT().ProductDiscounts(gomock.Any(), testProductID).
					Return(nil, errs.ErrNotFound).AnyTimes() // Используем AnyTimes() так как вызов может быть несколько раз

				// Mock order creation
				mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req dto.CreateOrderRepoReq) error {
						assert.Equal(t, testUserID, req.Order.UserID)
						assert.Equal(t, models.Placed, req.Order.Status)
						assert.Equal(t, 200.0, req.Order.TotalPrice)
						assert.Equal(t, 200.0, req.Order.TotalPriceDiscount)
						assert.Equal(t, testAddressID, req.Order.AddressID)
						assert.Len(t, req.Order.Items, 1)
						assert.Equal(t, testProductID, req.Order.Items[0].ProductID)
						assert.Equal(t, uint(2), req.Order.Items[0].Quantity)
						assert.Equal(t, 100.0, req.Order.Items[0].Price)
						assert.Equal(t, uint(3), req.UpdatedQuantities[testProductID])
						return nil
					})
			},
			expectedError: nil,
		},
		{
			name: "Product not approved",
			input: dto.CreateOrderDTO{
				UserID:    testUserID,
				AddressID: testAddressID,
				Items: []dto.CreateOrderItemDTO{
					{
						ProductID: testProductID,
						Quantity:  1,
					},
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().ProductPrice(gomock.Any(), testProductID).
					Return(&models.Product{
						ID:       testProductID,
						Status:   models.ProductPending,
						Quantity: 5,
					}, nil)
				// Добавляем ожидание вызова ProductDiscounts, даже если он не должен произойти
				mockRepo.EXPECT().ProductDiscounts(gomock.Any(), testProductID).
					Return(nil, errs.ErrNotFound).MaxTimes(0)
			},
			expectedError: errs.ErrProductNotApproved,
		},
		{
			name: "Not enough stock",
			input: dto.CreateOrderDTO{
				UserID:    testUserID,
				AddressID: testAddressID,
				Items: []dto.CreateOrderItemDTO{
					{
						ProductID: testProductID,
						Quantity:  10,
					},
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().ProductPrice(gomock.Any(), testProductID).
					Return(&models.Product{
						ID:       testProductID,
						Price:    100.0,
						Quantity: 5,
						Status:   models.ProductApproved,
					}, nil)
				// Добавляем ожидание вызова ProductDiscounts, даже если он не должен произойти
				mockRepo.EXPECT().ProductDiscounts(gomock.Any(), testProductID).
					Return(nil, errs.ErrNotFound).MaxTimes(0)
			},
			expectedError: errs.ErrNotEnoughStock,
		},
		{
			name: "Error getting product price",
			input: dto.CreateOrderDTO{
				UserID:    testUserID,
				AddressID: testAddressID,
				Items: []dto.CreateOrderItemDTO{
					{
						ProductID: testProductID,
						Quantity:  1,
					},
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().ProductPrice(gomock.Any(), testProductID).
					Return(nil, errors.New("database error"))
				// Добавляем ожидание вызова ProductDiscounts, даже если он не должен произойти
				mockRepo.EXPECT().ProductDiscounts(gomock.Any(), testProductID).
					Return(nil, errs.ErrNotFound).MaxTimes(0)
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := orderUC.CreateOrder(context.Background(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
