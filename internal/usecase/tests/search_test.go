package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/search"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchUsecase_SearchProductsByNameWithFilterAndSort(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISearchRepository(ctrl)
	uc := search.NewSearchUsecase(mockRepo)

	t.Run("success with sorting", func(t *testing.T) {
		ctx := context.Background()
		products := []*models.Product{
			{Name: "Product 1", Price: 100, PriceDiscount: 90, Rating: 4.5},
			{Name: "Product 2", Price: 200, PriceDiscount: 0, Rating: 3.5},
		}

		mockRepo.EXPECT().
			GetProductsByNameWithFilterAndSort(
				gomock.Any(),
				"test",
				null.String{},
				0,
				0.0,
				0.0,
				float32(0.0),
				models.SortByPriceAsc,
			).
			Return(products, nil)

		result, err := uc.SearchProductsByNameWithFilterAndSort(
			ctx,
			null.String{},
			"test",
			0,
			0.0,
			0.0,
			0.0,
			models.SortByPriceAsc,
		)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Product 1", result[0].Name)
		assert.Equal(t, "Product 2", result[1].Name)
	})

	t.Run("success with price desc sorting", func(t *testing.T) {
		ctx := context.Background()
		products := []*models.Product{
			{Name: "Product 1", Price: 100, PriceDiscount: 90, Rating: 4.5},
			{Name: "Product 2", Price: 200, PriceDiscount: 0, Rating: 3.5},
		}

		mockRepo.EXPECT().
			GetProductsByNameWithFilterAndSort(
				gomock.Any(),
				"test",
				null.String{},
				0,
				0.0,
				0.0,
				float32(0.0),
				models.SortByPriceDesc,
			).
			Return(products, nil)

		result, err := uc.SearchProductsByNameWithFilterAndSort(
			ctx,
			null.String{},
			"test",
			0,
			0.0,
			0.0,
			0.0,
			models.SortByPriceDesc,
		)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Product 2", result[0].Name)
		assert.Equal(t, "Product 1", result[1].Name)
	})

	t.Run("success with rating sorting", func(t *testing.T) {
		ctx := context.Background()
		products := []*models.Product{
			{Name: "Product 1", Price: 100, PriceDiscount: 90, Rating: 4.5},
			{Name: "Product 2", Price: 200, PriceDiscount: 0, Rating: 3.5},
		}

		mockRepo.EXPECT().
			GetProductsByNameWithFilterAndSort(
				gomock.Any(),
				"test",
				null.String{},
				0,
				0.0,
				0.0,
				float32(0.0),
				models.SortByRatingDesc,
			).
			Return(products, nil)

		result, err := uc.SearchProductsByNameWithFilterAndSort(
			ctx,
			null.String{},
			"test",
			0,
			0.0,
			0.0,
			0.0,
			models.SortByRatingDesc,
		)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Product 1", result[0].Name)
		assert.Equal(t, "Product 2", result[1].Name)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.Background()

		mockRepo.EXPECT().
			GetProductsByNameWithFilterAndSort(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).
			Return(nil, errors.New("repository error"))

		_, err := uc.SearchProductsByNameWithFilterAndSort(
			ctx,
			null.String{},
			"test",
			0,
			0.0,
			0.0,
			0.0,
			models.SortByDefault,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "repository error")
	})
}

func TestSearchUsecase_SearchCategoryByName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISearchRepository(ctrl)
	uc := search.NewSearchUsecase(mockRepo)

	t.Run("success multiple categories", func(t *testing.T) {
		ctx := context.Background()
		req := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{
				{Name: "Category 1"},
				{Name: "Category 2"},
			},
		}

		category1 := &models.Category{ID: uuid.New(), Name: "Category 1"}
		category2 := &models.Category{ID: uuid.New(), Name: "Category 2"}

		mockRepo.EXPECT().
			GetCategoryByName(gomock.Any(), "Category 1").
			Return(category1, nil)
		mockRepo.EXPECT().
			GetCategoryByName(gomock.Any(), "Category 2").
			Return(category2, nil)

		result, err := uc.SearchCategoryByName(ctx, req)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Contains(t, []string{result[0].Name, result[1].Name}, "Category 1")
		assert.Contains(t, []string{result[0].Name, result[1].Name}, "Category 2")
	})

	t.Run("some categories not found", func(t *testing.T) {
		ctx := context.Background()
		req := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{
				{Name: "Category 1"},
				{Name: "Non-existent"},
			},
		}

		category1 := &models.Category{ID: uuid.New(), Name: "Category 1"}

		mockRepo.EXPECT().
			GetCategoryByName(gomock.Any(), "Category 1").
			Return(category1, nil)
		mockRepo.EXPECT().
			GetCategoryByName(gomock.Any(), "Non-existent").
			Return(nil, nil)

		result, err := uc.SearchCategoryByName(ctx, req)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Category 1", result[0].Name)
	})

	t.Run("empty request", func(t *testing.T) {
		ctx := context.Background()
		req := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{},
		}

		result, err := uc.SearchCategoryByName(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("error in one goroutine", func(t *testing.T) {
		ctx := context.Background()
		req := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{
				{Name: "Category 1"},
				{Name: "Category 2"},
			},
		}

		mockRepo.EXPECT().
			GetCategoryByName(gomock.Any(), "Category 1").
			Return(nil, errors.New("database error"))
		mockRepo.EXPECT().
			GetCategoryByName(gomock.Any(), "Category 2").
			Return(&models.Category{Name: "Category 2"}, nil)

		_, err := uc.SearchCategoryByName(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		req := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{
				{Name: "Category 1"},
				{Name: "Category 2"},
			},
		}

		cancel()

		result, err := uc.SearchCategoryByName(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}