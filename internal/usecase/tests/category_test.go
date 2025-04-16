package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/category"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
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