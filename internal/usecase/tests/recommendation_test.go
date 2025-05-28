package tests

import (
	"context"
	"errors"
	"testing"

	recommendationRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	mockProduct "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/recommendation"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetRecommendations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRecommendationRepo := recommendationRepo.NewMockIRecommendationRepository(ctrl)
	mockProductUsecase := mockProduct.NewMockIProductUsecase(ctrl)

	usecase := recommendation.NewRecommendationUsecase(mockProductUsecase, mockRecommendationRepo)

	productID := uuid.New()
	subcatID1 := uuid.New()
	subcatID2 := uuid.New()

	// Product IDs returned from subcategories
	productIDs1 := []uuid.UUID{uuid.New(), uuid.New()}
	productIDs2 := []uuid.UUID{uuid.New(), productID} // includes original productID to check filtering

	expectedProductIDs := []uuid.UUID{productIDs1[0], productIDs1[1], productIDs2[0]} // original ID excluded

	expectedProducts := []*models.Product{
		{ID: productIDs1[0], Name: "Product 1"},
		{ID: productIDs1[1], Name: "Product 2"},
		{ID: productIDs2[0], Name: "Product 3"},
	}

	mockRecommendationRepo.EXPECT().
		GetCategoryIDsByProductID(gomock.Any(), productID).
		Return([]uuid.UUID{subcatID1, subcatID2}, nil)

	mockRecommendationRepo.EXPECT().
		GetProductIDsBySubcategoryID(gomock.Any(), subcatID1, 10).
		Return(productIDs1, nil)

	mockRecommendationRepo.EXPECT().
		GetProductIDsBySubcategoryID(gomock.Any(), subcatID2, 10).
		Return(productIDs2, nil)

	mockProductUsecase.EXPECT().
		GetProductsByIDs(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, ids []uuid.UUID) ([]*models.Product, error) {
			assert.ElementsMatch(t, expectedProductIDs, ids)
			return expectedProducts, nil
		})

	products, err := usecase.GetRecommendations(context.Background(), productID)
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
}

func TestGetRecommendations_GetCategoryIDsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRecommendationRepo := recommendationRepo.NewMockIRecommendationRepository(ctrl)
	mockProductUsecase := mockProduct.NewMockIProductUsecase(ctrl)

	usecase := recommendation.NewRecommendationUsecase(mockProductUsecase, mockRecommendationRepo)

	productID := uuid.New()

	mockRecommendationRepo.EXPECT().
		GetCategoryIDsByProductID(gomock.Any(), productID).
		Return(nil, errors.New("db error"))

	products, err := usecase.GetRecommendations(context.Background(), productID)
	assert.Error(t, err)
	assert.Nil(t, products)
}

func TestGetRecommendations_NoProductsFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRecommendationRepo := recommendationRepo.NewMockIRecommendationRepository(ctrl)
	mockProductUsecase := mockProduct.NewMockIProductUsecase(ctrl)

	usecase := recommendation.NewRecommendationUsecase(mockProductUsecase, mockRecommendationRepo)

	productID := uuid.New()
	subcatID := uuid.New()

	mockRecommendationRepo.EXPECT().
		GetCategoryIDsByProductID(gomock.Any(), productID).
		Return([]uuid.UUID{subcatID}, nil)

	mockRecommendationRepo.EXPECT().
		GetProductIDsBySubcategoryID(gomock.Any(), subcatID, 10).
		Return([]uuid.UUID{}, nil)

	products, err := usecase.GetRecommendations(context.Background(), productID)
	assert.NoError(t, err)
	assert.Nil(t, products)
}

func TestGetRecommendations_ProductUsecaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRecommendationRepo := recommendationRepo.NewMockIRecommendationRepository(ctrl)
	mockProductUsecase := mockProduct.NewMockIProductUsecase(ctrl)

	usecase := recommendation.NewRecommendationUsecase(mockProductUsecase, mockRecommendationRepo)

	productID := uuid.New()
	subcatID := uuid.New()
	productIDs := []uuid.UUID{uuid.New()}

	mockRecommendationRepo.EXPECT().
		GetCategoryIDsByProductID(gomock.Any(), productID).
		Return([]uuid.UUID{subcatID}, nil)

	mockRecommendationRepo.EXPECT().
		GetProductIDsBySubcategoryID(gomock.Any(), subcatID, 10).
		Return(productIDs, nil)

	mockProductUsecase.EXPECT().
		GetProductsByIDs(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("usecase error"))

	products, err := usecase.GetRecommendations(context.Background(), productID)
	assert.Error(t, err)
	assert.Nil(t, products)
}
