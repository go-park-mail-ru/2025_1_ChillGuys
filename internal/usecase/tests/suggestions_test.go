package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/golang/mock/gomock"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/suggestions"
)

func TestSuggestionsUsecase_GetCategorySuggestions(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISuggestionsRepository(ctrl)
	mockRedis := mocks.NewMockISuggestionsRedisRepository(ctrl)

	uc := suggestions.NewSuggestionsUsecase(mockRepo, mockRedis)

	t.Run("success from redis", func(t *testing.T) {
		ctx := context.Background()
		subString := "cat"
		redisNames := []string{"Category 1", "Category 2", "Catalog"}

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.CategoryNamesKey).
			Return(redisNames, nil)

		expected := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{
				{Name: "Category 1"},
				{Name: "Category 2"},
				{Name: "Catalog"},
			},
		}

		result, err := uc.GetCategorySuggestions(ctx, subString)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("success from db when redis empty", func(t *testing.T) {
		ctx := context.Background()
		subString := "cat"
		dbCategories := []*models.CategorySuggestion{
			{Name: "Category 1"},
			{Name: "Category 2"},
		}

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.CategoryNamesKey).
			Return([]string{}, nil)

		mockRepo.EXPECT().
			GetAllCategoriesName(ctx).
			Return(dbCategories, nil)

		mockRedis.EXPECT().
			AddSuggestionsByKey(ctx, redis.CategoryNamesKey, []string{"Category 1", "Category 2"}).
			Return(nil)

		expected := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{
				{Name: "Category 1"},
				{Name: "Category 2"},
			},
		}

		result, err := uc.GetCategorySuggestions(ctx, subString)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("redis error fallback to db", func(t *testing.T) {
		ctx := context.Background()
		subString := "cat"
		dbCategories := []*models.CategorySuggestion{
			{Name: "Category 1"},
		}

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.CategoryNamesKey).
			Return(nil, errors.New("redis error"))

		mockRepo.EXPECT().
			GetAllCategoriesName(ctx).
			Return(dbCategories, nil)

		mockRedis.EXPECT().
			AddSuggestionsByKey(ctx, redis.CategoryNamesKey, []string{"Category 1"}).
			Return(nil)

		expected := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{
				{Name: "Category 1"},
			},
		}

		result, err := uc.GetCategorySuggestions(ctx, subString)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("db error", func(t *testing.T) {
		ctx := context.Background()
		subString := "cat"

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.CategoryNamesKey).
			Return([]string{}, nil)

		mockRepo.EXPECT().
			GetAllCategoriesName(ctx).
			Return(nil, errors.New("db error"))

		_, err := uc.GetCategorySuggestions(ctx, subString)
		assert.Error(t, err)
	})

	t.Run("filtering works correctly", func(t *testing.T) {
		ctx := context.Background()
		subString := "elec"
		names := []string{"Electronics", "Books", "Electronic Gadgets"}

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.CategoryNamesKey).
			Return(names, nil)

		expected := dto.CategoryNameResponse{
			CategoriesNames: []models.CategorySuggestion{
				{Name: "Electronics"},
				{Name: "Electronic Gadgets"},
			},
		}

		result, err := uc.GetCategorySuggestions(ctx, subString)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestSuggestionsUsecase_GetProductSuggestions(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISuggestionsRepository(ctrl)
	mockRedis := mocks.NewMockISuggestionsRedisRepository(ctrl)

	uc := suggestions.NewSuggestionsUsecase(mockRepo, mockRedis)

	t.Run("success general from redis", func(t *testing.T) {
		ctx := context.Background()
		redisNames := []string{"Product 1", "Product 2"}

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.ProductNamesKey).
			Return(redisNames, nil)

		expected := dto.ProductNameResponse{
			ProductNames: []models.ProductSuggestion{
				{Name: "Product 1"},
				{Name: "Product 2"},
			},
		}

		result, err := uc.GetProductSuggestions(ctx, null.String{}, "")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("success by category from redis", func(t *testing.T) {
		ctx := context.Background()
		categoryID := "123e4567-e89b-12d3-a456-426614174000"
		redisNames := []string{"Category Product 1", "Category Product 2"}

		mockRedis.EXPECT().
			GetProductSuggestionsByCategory(ctx, categoryID).
			Return(redisNames, nil)

		expected := dto.ProductNameResponse{
			ProductNames: []models.ProductSuggestion{
				{Name: "Category Product 1"},
				{Name: "Category Product 2"},
			},
		}

		result, err := uc.GetProductSuggestions(ctx, null.StringFrom(categoryID), "")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("general from db when redis empty", func(t *testing.T) {
		ctx := context.Background()
		dbProducts := []*models.ProductSuggestion{
			{Name: "Product 1"},
			{Name: "Product 2"},
		}

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.ProductNamesKey).
			Return([]string{}, nil)

		mockRepo.EXPECT().
			GetAllProductsName(ctx).
			Return(dbProducts, nil)

		mockRedis.EXPECT().
			AddSuggestionsByKey(ctx, redis.ProductNamesKey, []string{"Product 1", "Product 2"}).
			Return(nil)

		expected := dto.ProductNameResponse{
			ProductNames: []models.ProductSuggestion{
				{Name: "Product 1"},
				{Name: "Product 2"},
			},
		}

		result, err := uc.GetProductSuggestions(ctx, null.String{}, "")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("by category from db when redis empty", func(t *testing.T) {
		ctx := context.Background()
		categoryID := "123e4567-e89b-12d3-a456-426614174000"
		dbProducts := []*models.ProductSuggestion{
			{Name: "Category Product 1"},
		}

		mockRedis.EXPECT().
			GetProductSuggestionsByCategory(ctx, categoryID).
			Return([]string{}, nil)

		mockRepo.EXPECT().
			GetProductsNameByCategory(ctx, categoryID).
			Return(dbProducts, nil)

		mockRedis.EXPECT().
			AddProductSuggestionsByCategory(ctx, categoryID, []string{"Category Product 1"}).
			Return(nil)

		expected := dto.ProductNameResponse{
			ProductNames: []models.ProductSuggestion{
				{Name: "Category Product 1"},
			},
		}

		result, err := uc.GetProductSuggestions(ctx, null.StringFrom(categoryID), "")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("redis error fallback to db for general", func(t *testing.T) {
		ctx := context.Background()
		dbProducts := []*models.ProductSuggestion{
			{Name: "Product 1"},
		}

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.ProductNamesKey).
			Return(nil, errors.New("redis error"))

		mockRepo.EXPECT().
			GetAllProductsName(ctx).
			Return(dbProducts, nil)

		mockRedis.EXPECT().
			AddSuggestionsByKey(ctx, redis.ProductNamesKey, []string{"Product 1"}).
			Return(nil)

		expected := dto.ProductNameResponse{
			ProductNames: []models.ProductSuggestion{
				{Name: "Product 1"},
			},
		}

		result, err := uc.GetProductSuggestions(ctx, null.String{}, "")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("db error for general", func(t *testing.T) {
		ctx := context.Background()

		mockRedis.EXPECT().
			GetSuggestionsByKey(ctx, redis.ProductNamesKey).
			Return([]string{}, nil)

		mockRepo.EXPECT().
			GetAllProductsName(ctx).
			Return(nil, errors.New("db error"))

		_, err := uc.GetProductSuggestions(ctx, null.String{}, "")
		assert.Error(t, err)
	})
}