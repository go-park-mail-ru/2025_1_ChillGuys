package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/category"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategoryUsecase_GetAllCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICategoryRepository(ctrl)
	uc := category.NewCategoryUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedCategories := []*models.Category{
			{
				ID:   uuid.New(),
				Name: "Category 1",
			},
			{
				ID:   uuid.New(),
				Name: "Category 2",
			},
		}

		mockRepo.EXPECT().
			GetAllCategories(gomock.Any()).
			Return(expectedCategories, nil)

		categories, err := uc.GetAllCategories(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
	})

	t.Run("empty result", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllCategories(gomock.Any()).
			Return([]*models.Category{}, nil)

		categories, err := uc.GetAllCategories(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, categories)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllCategories(gomock.Any()).
			Return(nil, errors.New("database error"))

		categories, err := uc.GetAllCategories(context.Background())
		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Contains(t, err.Error(), "CategoryUsecase.GetAllCategories")
	})
}

func TestCategoryUsecase_GetAllSubcategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICategoryRepository(ctrl)
	uc := category.NewCategoryUsecase(mockRepo)

	categoryID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedSubcategories := []*models.Category{
			{
				ID:   uuid.New(),
				Name: "Subcategory 1",
			},
			{
				ID:   uuid.New(),
				Name: "Subcategory 2",
			},
		}

		mockRepo.EXPECT().
			GetAllSubcategories(gomock.Any(), categoryID).
			Return(expectedSubcategories, nil)

		subcategories, err := uc.GetAllSubategories(context.Background(), categoryID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSubcategories, subcategories)
	})

	t.Run("empty result", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllSubcategories(gomock.Any(), categoryID).
			Return([]*models.Category{}, nil)

		subcategories, err := uc.GetAllSubategories(context.Background(), categoryID)
		assert.NoError(t, err)
		assert.Empty(t, subcategories)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllSubcategories(gomock.Any(), categoryID).
			Return(nil, errors.New("database error"))

		subcategories, err := uc.GetAllSubategories(context.Background(), categoryID)
		assert.Error(t, err)
		assert.Nil(t, subcategories)
		assert.Contains(t, err.Error(), "CategoryUsecase.GetAllCategories") // Note: This matches your implementation's op name
	})
}

func TestCategoryUsecase_GetNameSubcategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICategoryRepository(ctrl)
	uc := category.NewCategoryUsecase(mockRepo)

	subcategoryID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedName := "Test Subcategory"

		mockRepo.EXPECT().
			GetNameSubcategory(gomock.Any(), subcategoryID).
			Return(expectedName, nil)

		name, err := uc.GetNameSubcategory(context.Background(), subcategoryID)
		assert.NoError(t, err)
		assert.Equal(t, expectedName, name)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().
			GetNameSubcategory(gomock.Any(), subcategoryID).
			Return("", errors.New("database error"))

		name, err := uc.GetNameSubcategory(context.Background(), subcategoryID)
		assert.Error(t, err)
		assert.Empty(t, name)
		assert.Contains(t, err.Error(), "CategoryUsecase.GetNameSubcategory")
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetNameSubcategory(gomock.Any(), subcategoryID).
			Return("", nil) // Assuming empty string means not found

		name, err := uc.GetNameSubcategory(context.Background(), subcategoryID)
		assert.NoError(t, err) // Depending on your business logic, you might want to return an error here
		assert.Empty(t, name)
	})
}
