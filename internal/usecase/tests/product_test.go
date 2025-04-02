package tests_test

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

func TestProductUsecase_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	logger := logrus.New()
	uc := product.NewProductUsecase(logger, mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedProducts := []*models.Product{
			{
				ID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
				Name: "Product 1",
			},
			{
				ID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
				Name: "Product 2",
			},
		}

		mockRepo.EXPECT().
			GetAllProducts(gomock.Any()).
			Return(expectedProducts, nil).
			Times(1)

		products, err := uc.GetAllProducts(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllProducts(gomock.Any()).
			Return(nil, errors.New("infrastructure error")).
			Times(1)

		products, err := uc.GetAllProducts(context.Background())

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "infrastructure error")
	})
}

func TestProductUsecase_GetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	logger := logrus.New()
	uc := product.NewProductUsecase(logger, mockRepo)

	t.Run("Success", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		expectedProduct := &models.Product{
			ID:          testID,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       1000,
			Quantity:    5,
		}

		mockRepo.EXPECT().
			GetProductByID(gomock.Any(), testID).
			Return(expectedProduct, nil).
			Times(1)

		product, err := uc.GetProductByID(context.Background(), testID)

		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("NotFound", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440999")

		mockRepo.EXPECT().
			GetProductByID(gomock.Any(), testID).
			Return(nil, errors.New("not found")).
			Times(1)

		product, err := uc.GetProductByID(context.Background(), testID)

		assert.Error(t, err)
		assert.Nil(t, product)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

		mockRepo.EXPECT().
			GetProductByID(gomock.Any(), testID).
			Return(nil, errors.New("database error")).
			Times(1)

		product, err := uc.GetProductByID(context.Background(), testID)

		assert.Error(t, err)
		assert.Nil(t, product)
	})
}

func TestProductUsecase_GetProductCover(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIProductRepository(ctrl)
	logger := logrus.New()
	uc := product.NewProductUsecase(logger, mockRepo)

	t.Run("Success", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		expectedData := []byte{0xFF, 0xD8, 0xFF} // Пример JPEG данных

		mockRepo.EXPECT().
			GetProductCoverPath(gomock.Any(), testID).
			Return(expectedData, nil).
			Times(1)

		data, err := uc.GetProductCover(context.Background(), testID)

		assert.NoError(t, err)
		assert.Equal(t, expectedData, data)
	})

	t.Run("NotFound", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

		mockRepo.EXPECT().
			GetProductCoverPath(gomock.Any(), testID).
			Return(nil, errors.New("file not found")).
			Times(1)

		data, err := uc.GetProductCover(context.Background(), testID)

		assert.Error(t, err)
		assert.Nil(t, data)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")

		mockRepo.EXPECT().
			GetProductCoverPath(gomock.Any(), testID).
			Return(nil, errors.New("storage error")).
			Times(1)

		data, err := uc.GetProductCover(context.Background(), testID)

		assert.Error(t, err)
		assert.Nil(t, data)
	})
}
