package tests

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
// 	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
// 	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
// 	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
// 	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/guregu/null"
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// )

// func TestCreateOrder(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockIOrderRepository(ctrl)

// 	orderUC := order.NewOrderUsecase(mockRepo)

// 	testUserID := uuid.New()
// 	testAddressID := uuid.New()
// 	testProductID := uuid.New()

// 	tests := []struct {
// 		name          string
// 		input         dto.CreateOrderDTO
// 		mockSetup     func()
// 		expectedError error
// 	}{
// 		{
// 			name: "Successful order creation",
// 			input: dto.CreateOrderDTO{
// 				UserID:    testUserID,
// 				AddressID: testAddressID,
// 				Items: []dto.CreateOrderItemDTO{
// 					{
// 						ProductID: testProductID,
// 						Quantity:  2,
// 					},
// 				},
// 			},
// 			mockSetup: func() {
// 				// Mock product price check
// 				mockRepo.EXPECT().ProductPrice(gomock.Any(), testProductID).
// 					Return(&models.Product{
// 						ID:       testProductID,
// 						Price:    100.0,
// 						Quantity: 5,
// 						Status:   models.ProductApproved,
// 					}, nil)

// 				// Mock discount check (no discount)
// 				mockRepo.EXPECT().ProductDiscounts(gomock.Any(), testProductID).
// 					Return(nil, errs.ErrNotFound).AnyTimes() // Используем AnyTimes() так как вызов может быть несколько раз

// 				// Mock order creation
// 				mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).
// 					DoAndReturn(func(ctx context.Context, req dto.CreateOrderRepoReq) error {
// 						assert.Equal(t, testUserID, req.Order.UserID)
// 						assert.Equal(t, models.Placed, req.Order.Status)
// 						assert.Equal(t, 200.0, req.Order.TotalPrice)
// 						assert.Equal(t, 200.0, req.Order.TotalPriceDiscount)
// 						assert.Equal(t, testAddressID, req.Order.AddressID)
// 						assert.Len(t, req.Order.Items, 1)
// 						assert.Equal(t, testProductID, req.Order.Items[0].ProductID)
// 						assert.Equal(t, uint(2), req.Order.Items[0].Quantity)
// 						assert.Equal(t, 100.0, req.Order.Items[0].Price)
// 						assert.Equal(t, uint(3), req.UpdatedQuantities[testProductID])
// 						return nil
// 					})
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Product not approved",
// 			input: dto.CreateOrderDTO{
// 				UserID:    testUserID,
// 				AddressID: testAddressID,
// 				Items: []dto.CreateOrderItemDTO{
// 					{
// 						ProductID: testProductID,
// 						Quantity:  1,
// 					},
// 				},
// 			},
// 			mockSetup: func() {
// 				mockRepo.EXPECT().ProductPrice(gomock.Any(), testProductID).
// 					Return(&models.Product{
// 						ID:       testProductID,
// 						Status:   models.ProductPending,
// 						Quantity: 5,
// 					}, nil)
// 				// Добавляем ожидание вызова ProductDiscounts, даже если он не должен произойти
// 				mockRepo.EXPECT().ProductDiscounts(gomock.Any(), testProductID).
// 					Return(nil, errs.ErrNotFound).MaxTimes(0)
// 			},
// 			expectedError: errs.ErrProductNotApproved,
// 		},
// 		{
// 			name: "Not enough stock",
// 			input: dto.CreateOrderDTO{
// 				UserID:    testUserID,
// 				AddressID: testAddressID,
// 				Items: []dto.CreateOrderItemDTO{
// 					{
// 						ProductID: testProductID,
// 						Quantity:  10,
// 					},
// 				},
// 			},
// 			mockSetup: func() {
// 				mockRepo.EXPECT().ProductPrice(gomock.Any(), testProductID).
// 					Return(&models.Product{
// 						ID:       testProductID,
// 						Price:    100.0,
// 						Quantity: 5,
// 						Status:   models.ProductApproved,
// 					}, nil)
// 				// Добавляем ожидание вызова ProductDiscounts, даже если он не должен произойти
// 				mockRepo.EXPECT().ProductDiscounts(gomock.Any(), testProductID).
// 					Return(nil, errs.ErrNotFound).MaxTimes(0)
// 			},
// 			expectedError: errs.ErrNotEnoughStock,
// 		},
// 		{
// 			name: "Error getting product price",
// 			input: dto.CreateOrderDTO{
// 				UserID:    testUserID,
// 				AddressID: testAddressID,
// 				Items: []dto.CreateOrderItemDTO{
// 					{
// 						ProductID: testProductID,
// 						Quantity:  1,
// 					},
// 				},
// 			},
// 			mockSetup: func() {
// 				mockRepo.EXPECT().ProductPrice(gomock.Any(), testProductID).
// 					Return(nil, errors.New("database error"))
// 				// Добавляем ожидание вызова ProductDiscounts, даже если он не должен произойти
// 				mockRepo.EXPECT().ProductDiscounts(gomock.Any(), testProductID).
// 					Return(nil, errs.ErrNotFound).MaxTimes(0)
// 			},
// 			expectedError: errors.New("database error"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSetup()

// 			err := orderUC.CreateOrder(context.Background(), tt.input)

// 			if tt.expectedError != nil {
// 				assert.Error(t, err)
// 				assert.Contains(t, err.Error(), tt.expectedError.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestGetUserOrders(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockIOrderRepository(ctrl)
// 	orderUC := order.NewOrderUsecase(mockRepo)

// 	testUserID := uuid.New()
// 	orderID := uuid.New()
// 	addressID := uuid.New()
// 	productID := uuid.New()

// 	tests := []struct {
// 		name           string
// 		mockSetup      func()
// 		expectedOrders *[]dto.OrderPreviewDTO
// 		expectedError  error
// 	}{
// 		{
// 			name: "Successful retrieval of orders",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().GetOrdersByUserID(gomock.Any(), testUserID).
// 					Return(&[]dto.GetOrderByUserIDResDTO{
// 						{
// 							ID:                 orderID,
// 							Status:             models.Placed,
// 							TotalPrice:         100.0,
// 							TotalPriceDiscount: 10.0,
// 							AddressID:          addressID,
// 							ExpectedDeliveryAt: nil,
// 							ActualDeliveryAt:   nil,
// 							CreatedAt:          nil,
// 						},
// 					}, nil)

// 				mockRepo.EXPECT().GetOrderProducts(gomock.Any(), orderID).
// 					Return(&[]dto.GetOrderProductResDTO{
// 						{
// 							ProductID: productID,
// 							Quantity:  2,
// 						},
// 					}, nil)

// 				mockRepo.EXPECT().GetOrderAddress(gomock.Any(), addressID).
// 					Return(&models.AddressDB{ID: addressID, City: null.StringFrom("City")}, nil)

// 				mockRepo.EXPECT().GetProductImage(gomock.Any(), productID).
// 					Return("https://example.com/image.jpg", nil)
// 			},
// 			expectedOrders: &[]dto.OrderPreviewDTO{
// 				{
// 					ID:                 orderID,
// 					Status:             models.Placed,
// 					TotalPrice:         100.0,
// 					TotalDiscountPrice: 10.0,
// 					Products: []models.OrderPreviewProductDTO{
// 						{
// 							ProductImageURL: null.StringFrom("https://example.com/image.jpg"),
// 							ProductQuantity: 2,
// 						},
// 					},
// 					Address:            models.AddressDB{ID: addressID, City: null.StringFrom("City")},
// 					ExpectedDeliveryAt: nil,
// 					ActualDeliveryAt:   nil,
// 					CreatedAt:          nil,
// 				},
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Error retrieving orders",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().GetOrdersByUserID(gomock.Any(), testUserID).
// 					Return(nil, errors.New("database error"))
// 			},
// 			expectedOrders: nil,
// 			expectedError:  fmt.Errorf("OrderUsecase.GetUserOrders: database error"),
// 		},
// 		{
// 			name: "No orders found",
// 			mockSetup: func() {
// 				mockRepo.EXPECT().GetOrdersByUserID(gomock.Any(), testUserID).
// 					Return(nil, errs.ErrNotFound)
// 			},
// 			expectedOrders: nil,
// 			expectedError:  fmt.Errorf("OrderUsecase.GetUserOrders: %w", errs.NewNotFoundError("OrderUsecase.GetUserOrders")),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSetup()

// 			orders, err := orderUC.GetUserOrders(context.Background(), testUserID)

// 			if tt.expectedError != nil {
// 				assert.Error(t, err)
// 				assert.Contains(t, err.Error(), tt.expectedError.Error())
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, *tt.expectedOrders, *orders)
// 			}
// 		})
// 	}
// }
