package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductUsecase_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	uc := product.NewProductUsecase(mockRepo)

	tests := []struct {
		name          string
		offset        int
		mockSetup     func()
		expected      []*models.Product
		expectedError error
	}{
		{
			name:   "success",
			offset: 0,
			mockSetup: func() {
				expectedProducts := []*models.Product{
					{ID: uuid.New(), Name: "Product 1"},
					{ID: uuid.New(), Name: "Product 2"},
				}
				mockRepo.EXPECT().
					GetAllProducts(gomock.Any(), 0).
					Return(expectedProducts, nil)
			},
			expected: []*models.Product{
				{ID: uuid.New(), Name: "Product 1"},
				{ID: uuid.New(), Name: "Product 2"},
			},
		},
		{
			name:   "repository error",
			offset: 0,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetAllProducts(gomock.Any(), 0).
					Return(nil, errors.New("repository error"))
			},
			expectedError: errors.New("ProductUsecase.GetAllProducts: repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			products, err := uc.GetAllProducts(context.Background(), tt.offset)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(products))
		})
	}
}

func TestProductUsecase_GetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	uc := product.NewProductUsecase(mockRepo)

	tests := []struct {
		name          string
		productID     uuid.UUID
		mockSetup     func()
		expected      *models.Product
		expectedError error
	}{
		{
			name:      "success",
			productID: uuid.New(),
			mockSetup: func() {
				expectedProduct := &models.Product{
					ID:   uuid.New(),
					Name: "Test Product",
				}
				mockRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Any()).
					Return(expectedProduct, nil)
			},
			expected: &models.Product{
				ID:   uuid.New(),
				Name: "Test Product",
			},
		},
		{
			name:      "not found",
			productID: uuid.New(),
			mockSetup: func() {
				mockRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Any()).
					Return(nil, errs.ErrNotFound)
			},
			expectedError: errors.New("ProductUsecase.GetProductByID: not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			product, err := uc.GetProductByID(context.Background(), tt.productID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, product)
		})
	}
}

func TestProductUsecase_GetProductsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	uc := product.NewProductUsecase(mockRepo)

	tests := []struct {
		name          string
		categoryID    uuid.UUID
		offset        int
		minPrice      float64
		maxPrice      float64
		minRating     float32
		sortOption    models.SortOption
		mockSetup     func()
		expected      []*models.Product
		expectedError error
	}{
		{
			name:       "success with default sort",
			categoryID: uuid.New(),
			offset:     0,
			minPrice:   0,
			maxPrice:   0,
			minRating:  0,
			sortOption: models.SortByDefault,
			mockSetup: func() {
				expectedProducts := []*models.Product{
					{ID: uuid.New(), Name: "Product 1", Price: 1000, Rating: 4},
					{ID: uuid.New(), Name: "Product 2", Price: 2000, Rating: 5},
				}
				mockRepo.EXPECT().
					GetProductsByCategory(gomock.Any(), gomock.Any(), 0, 0.0, 0.0, float32(0), models.SortByDefault).
					Return(expectedProducts, nil)
			},
			expected: []*models.Product{
				{ID: uuid.New(), Name: "Product 1", Price: 1000, Rating: 4},
				{ID: uuid.New(), Name: "Product 2", Price: 2000, Rating: 5},
			},
		},
		{
			name:       "success with price asc sort",
			categoryID: uuid.New(),
			offset:     0,
			minPrice:   0,
			maxPrice:   0,
			minRating:  0,
			sortOption: models.SortByPriceAsc,
			mockSetup: func() {
				expectedProducts := []*models.Product{
					{ID: uuid.New(), Name: "Product 1", Price: 2000, Rating: 4},
					{ID: uuid.New(), Name: "Product 2", Price: 1000, Rating: 5},
				}
				mockRepo.EXPECT().
					GetProductsByCategory(gomock.Any(), gomock.Any(), 0, 0.0, 0.0, float32(0), models.SortByPriceAsc).
					Return(expectedProducts, nil)
			},
			expected: []*models.Product{
				{ID: uuid.New(), Name: "Product 2", Price: 1000, Rating: 5},
				{ID: uuid.New(), Name: "Product 1", Price: 2000, Rating: 4},
			},
		},
		{
			name:       "repository error",
			categoryID: uuid.New(),
			mockSetup: func() {
				mockRepo.EXPECT().
					GetProductsByCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
			expectedError: errors.New("ProductUsecase.GetProductsByCategoryWithFilterAndSort: repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			products, err := uc.GetProductsByCategory(
				context.Background(),
				tt.categoryID,
				tt.offset,
				tt.minPrice,
				tt.maxPrice,
				tt.minRating,
				tt.sortOption,
			)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(products))
		})
	}
}

func TestProductUsecase_GetProductsByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	uc := product.NewProductUsecase(mockRepo)

	tests := []struct {
		name          string
		ids           []uuid.UUID
		mockSetup     func()
		expected      []*models.Product
		expectedError error
	}{
		{
			name: "success",
			ids:  []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Any()).
					Return(&models.Product{ID: uuid.New()}, nil).
					Times(2)
			},
			expected: []*models.Product{
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
		},
		{
			name: "empty ids",
			ids:  []uuid.UUID{},
			mockSetup: func() {
				// No expectations - shouldn't call repository
			},
			expected: []*models.Product{},
		},
		{
			name: "repository error",
			ids:  []uuid.UUID{uuid.New()},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetProductByID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			products, err := uc.GetProductsByIDs(context.Background(), tt.ids)
			if tt.expectedError != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.ids), len(products))
		})
	}
}

func TestProductUsecase_AddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	uc := product.NewProductUsecase(mockRepo)

	tests := []struct {
		name          string
		product       *models.Product
		categoryID    uuid.UUID
		mockSetup     func()
		expected      *models.Product
		expectedError error
	}{
		{
			name: "success",
			product: &models.Product{
				Name:  "Test Product",
				Price: 1000,
			},
			categoryID: uuid.New(),
			mockSetup: func() {
				mockRepo.EXPECT().
					AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Product{
						ID:    uuid.New(),
						Name:  "Test Product",
						Price: 1000,
					}, nil)
			},
			expected: &models.Product{
				ID:    uuid.New(),
				Name:  "Test Product",
				Price: 1000,
			},
		},
		{
			name: "empty name",
			product: &models.Product{
				Name:  "",
				Price: 1000,
			},
			expectedError: errs.ErrEmptyProductName,
		},
		{
			name: "invalid price",
			product: &models.Product{
				Name:  "Test",
				Price: -100,
			},
			expectedError: errs.ErrInvalidProductPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			product, err := uc.AddProduct(context.Background(), tt.product, tt.categoryID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, product)
		})
	}
}