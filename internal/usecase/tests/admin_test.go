package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	admin "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/admin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAdminUsecase_GetPendingProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdminRepository(ctrl)
	mockRedisRepo := &redis.SuggestionsRepository{}
	mockProductRepo := mocks.NewMockIProductRepository(ctrl)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	sellerID := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()

	tests := []struct {
		name          string
		offset        int
		mockProducts  []*models.Product
		mockError     error
		expectedResp  dto.ProductsResponse
		expectedError error
	}{
		{
			name:   "Success",
			offset: 0,
			mockProducts: []*models.Product{
				{
					ID:              productID1,
					SellerID:        sellerID,
					Name:            "Product 1",
					PreviewImageURL: "image1.jpg",
					Description:     "Description 1",
					Status:          models.ProductPending,
					Price:           100.0,
					PriceDiscount:   90.0,
					Quantity:        10,
					Rating:          4.5,
					ReviewsCount:    20,
					UpdatedAt:       time.Now(),
				},
				{
					ID:              productID2,
					SellerID:        sellerID,
					Name:            "Product 2",
					PreviewImageURL: "image2.jpg",
					Description:     "Description 2",
					Status:          models.ProductPending,
					Price:           200.0,
					PriceDiscount:   180.0,
					Quantity:        20,
					Rating:          4.8,
					ReviewsCount:    30,
					UpdatedAt:       time.Now(),
				},
			},
			mockError: nil,
			expectedResp: dto.ProductsResponse{
				Total: 2,
				Products: []dto.BriefProduct{
					{
						ID:            productID1,
						Name:          "Product 1",
						ImageURL:      "image1.jpg",
						Price:         100.0,
						PriceDiscount: 90.0,
						Quantity:      10,
						Rating:        4.5,
						ReviewsCount:  20,
					},
					{
						ID:            productID2,
						Name:          "Product 2",
						ImageURL:      "image2.jpg",
						Price:         200.0,
						PriceDiscount: 180.0,
						Quantity:      20,
						Rating:        4.8,
						ReviewsCount:  30,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:          "Repository Error",
			offset:        0,
			mockProducts:  nil,
			mockError:     errors.New("repository error"),
			expectedResp:  dto.ProductsResponse{},
			expectedError: errors.New("AdminUsecase.GetPendingProducts: repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetPendingProducts(ctx, tt.offset).Return(tt.mockProducts, tt.mockError)

			uc := admin.NewAdminUsecase(mockRepo, mockRedisRepo, mockProductRepo)
			resp, err := uc.GetPendingProducts(ctx, tt.offset)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp.Total, resp.Total)
				assert.Equal(t, len(tt.expectedResp.Products), len(resp.Products))
				for i := range resp.Products {
					assert.Equal(t, tt.expectedResp.Products[i].ID, resp.Products[i].ID)
					assert.Equal(t, tt.expectedResp.Products[i].Name, resp.Products[i].Name)
					assert.Equal(t, tt.expectedResp.Products[i].ImageURL, resp.Products[i].ImageURL)
					assert.Equal(t, tt.expectedResp.Products[i].Price, resp.Products[i].Price)
					assert.Equal(t, tt.expectedResp.Products[i].PriceDiscount, resp.Products[i].PriceDiscount)
					assert.Equal(t, tt.expectedResp.Products[i].Quantity, resp.Products[i].Quantity)
					assert.Equal(t, tt.expectedResp.Products[i].Rating, resp.Products[i].Rating)
					assert.Equal(t, tt.expectedResp.Products[i].ReviewsCount, resp.Products[i].ReviewsCount)
				}
			}
		})
	}
}

func TestAdminUsecase_GetPendingUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdminRepository(ctrl)
	mockRedisRepo := &redis.SuggestionsRepository{}
	mockProductRepo := mocks.NewMockIProductRepository(ctrl)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	userID1 := uuid.New()
	userID2 := uuid.New()

	tests := []struct {
		name         string
		offset       int
		mockUsers    []*models.User
		mockError    error
		expectedResp dto.UsersResponse
		expectedErr  error
	}{
		{
			name:   "Success",
			offset: 0,
			mockUsers: []*models.User{
				{ID: userID1, Email: "user1@example.com", Role: models.RoleBuyer},
				{ID: userID2, Email: "user2@example.com", Role: models.RoleBuyer},
			},
			mockError: nil,
			expectedResp: dto.UsersResponse{
				Users: []dto.BriefUser{
					{ID: userID1, Email: "user1@example.com", Role: "buyer"},
					{ID: userID2, Email: "user2@example.com", Role: "buyer"},
				},
			},
			expectedErr: nil,
		},
		{
			name:         "Repository Error",
			offset:       0,
			mockUsers:    nil,
			mockError:    errors.New("repository error"),
			expectedResp: dto.UsersResponse{},
			expectedErr:  errors.New("AdminUsecase.GetPendingUsers: repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetPendingUsers(ctx, tt.offset).Return(tt.mockUsers, tt.mockError)

			uc := admin.NewAdminUsecase(mockRepo, mockRedisRepo, mockProductRepo)
			resp, err := uc.GetPendingUsers(ctx, tt.offset)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResp.Users), len(resp.Users))
				for i := range resp.Users {
					assert.Equal(t, tt.expectedResp.Users[i].ID, resp.Users[i].ID)
					assert.Equal(t, tt.expectedResp.Users[i].Email, resp.Users[i].Email)
					assert.Equal(t, tt.expectedResp.Users[i].Role, resp.Users[i].Role)
				}
			}
		})
	}
}

func TestAdminUsecase_UpdateUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdminRepository(ctrl)
	mockRedisRepo := &redis.SuggestionsRepository{}
	mockProductRepo := mocks.NewMockIProductRepository(ctrl)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	userID := uuid.New()

	tests := []struct {
		name          string
		req           dto.UpdateUserRoleRequest
		mockError     error
		expectedError error
	}{
		{
			name: "Promote to Seller Success",
			req: dto.UpdateUserRoleRequest{
				UserID: userID,
				Update: 1,
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "Demote to Buyer Success",
			req: dto.UpdateUserRoleRequest{
				UserID: userID,
				Update: 0,
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "Invalid Update Value",
			req: dto.UpdateUserRoleRequest{
				UserID: userID,
				Update: 2,
			},
			expectedError: errs.ErrParseRequestData,
		},
		{
			name: "Repository Error",
			req: dto.UpdateUserRoleRequest{
				UserID: userID,
				Update: 1,
			},
			mockError:     errors.New("repository error"),
			expectedError: errors.New("AdminUsecase.UpdateUserRole: repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Update == 0 || tt.req.Update == 1 {
				var expectedRole models.UserRole
				if tt.req.Update == 0 {
					expectedRole = models.RoleBuyer
				} else {
					expectedRole = models.RoleSeller
				}

				mockRepo.EXPECT().UpdateUserRole(ctx, tt.req.UserID, expectedRole).Return(tt.mockError)
			}

			uc := admin.NewAdminUsecase(mockRepo, mockRedisRepo, mockProductRepo)
			err := uc.UpdateUserRole(ctx, tt.req)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdminUsecase_UpdateProductStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdminRepository(ctrl)
	mockProductRepo := mocks.NewMockIProductRepository(ctrl)
	mockRedisRepo := &redis.SuggestionsRepository{} // конкретная реализация, как в твоем успешном примере

	uc := admin.NewAdminUsecase(mockRepo, mockRedisRepo, mockProductRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	productID := uuid.New() // UUID вместо int64

	tests := []struct {
		name          string
		updateValue   int
		setupMocks    func()
		expectedError error
	}{
		{
			name:          "некорректное значение статуса",
			updateValue:   999,
			setupMocks:    func() {},
			expectedError: errs.ErrParseRequestData,
		},
		{
			name:        "ошибка обновления статуса в БД",
			updateValue: 1,
			setupMocks: func() {
				mockRepo.EXPECT().
					UpdateProductStatus(ctx, productID, models.ProductApproved).
					Return(errors.New("db error"))
			},
			expectedError: errors.New("AdminUsecase.UpdateProductStatus: db error"),
		},
		{
			name:        "ошибка получения продукта",
			updateValue: 1,
			setupMocks: func() {
				mockRepo.EXPECT().
					UpdateProductStatus(ctx, productID, models.ProductApproved).
					Return(nil)

				mockProductRepo.EXPECT().
					GetProductByID(ctx, productID).
					Return(nil, errors.New("not found"))
			},
			expectedError: errors.New("AdminUsecase.UpdateProductStatus: not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := dto.UpdateProductStatusRequest{
				ProductID: productID,
				Update:    tt.updateValue,
			}

			tt.setupMocks()

			err := uc.UpdateProductStatus(ctx, req)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
