package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductUsecase_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	uc := product.NewProductUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedProducts := []*models.Product{
			{
				ID:   uuid.New(),
				Name: "Product 1",
			},
			{
				ID:   uuid.New(),
				Name: "Product 2",
			},
		}

		mockRepo.EXPECT().
			GetAllProducts(gomock.Any()).
			Return(expectedProducts, nil)

		products, err := uc.GetAllProducts(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllProducts(gomock.Any()).
			Return(nil, errors.New("repository error"))

		products, err := uc.GetAllProducts(context.Background())
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "ProductUsecase.GetAllProducts")
	})
}

func TestProductUsecase_GetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	uc := product.NewProductUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		productID := uuid.New()
		expectedProduct := &models.Product{
			ID:   productID,
			Name: "Test Product",
		}

		mockRepo.EXPECT().
			GetProductByID(gomock.Any(), productID).
			Return(expectedProduct, nil)

		product, err := uc.GetProductByID(context.Background(), productID)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("not found", func(t *testing.T) {
		productID := uuid.New()

		mockRepo.EXPECT().
			GetProductByID(gomock.Any(), productID).
			Return(nil, errors.New("not found"))

		product, err := uc.GetProductByID(context.Background(), productID)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "ProductUsecase.GetProductByID")
	})
}

func TestProductUsecase_GetProductsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	uc := product.NewProductUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		categoryID := uuid.New()
		expectedProducts := []*models.Product{
			{
				ID:   uuid.New(),
				Name: "Product 1",
			},
			{
				ID:   uuid.New(),
				Name: "Product 2",
			},
		}

		mockRepo.EXPECT().
			GetProductsByCategory(gomock.Any(), categoryID).
			Return(expectedProducts, nil)

		products, err := uc.GetProductsByCategory(context.Background(), categoryID)
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("empty result", func(t *testing.T) {
		categoryID := uuid.New()

		mockRepo.EXPECT().
			GetProductsByCategory(gomock.Any(), categoryID).
			Return([]*models.Product{}, nil)

		products, err := uc.GetProductsByCategory(context.Background(), categoryID)
		assert.NoError(t, err)
		assert.Empty(t, products)
	})

	t.Run("repository error", func(t *testing.T) {
		categoryID := uuid.New()

		mockRepo.EXPECT().
			GetProductsByCategory(gomock.Any(), categoryID).
			Return(nil, errors.New("repository error"))

		products, err := uc.GetProductsByCategory(context.Background(), categoryID)
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "ProductUsecase.GetProductsByCategory")
	})
}