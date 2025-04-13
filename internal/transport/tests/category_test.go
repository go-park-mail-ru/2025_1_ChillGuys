package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/category"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestCategory(t *testing.T) (*mocks.MockICategoryUsecase, *category.CategoryService) {
	ctrl := gomock.NewController(t)
	mockUsecase := mocks.NewMockICategoryUsecase(ctrl)
	service := category.NewCategoryService(mockUsecase)
	return mockUsecase, service
}

func TestCategoryService_GetAllCategories(t *testing.T) {
	mockUsecase, service := setupTestCategory(t)

	t.Run("success with categories", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()
		expectedCategories := []*models.Category{
			{
				ID:   id1,
				Name: "Category 1",
			},
			{
				ID:   id2,
				Name: "Category 2",
			},
		}

		mockUsecase.EXPECT().
			GetAllCategories(gomock.Any()).
			Return(expectedCategories, nil)

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		service.GetAllCategories(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseData dto.CategoryResponse
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)

		assert.Equal(t, len(expectedCategories), responseData.Total)
		assert.Equal(t, expectedCategories[0].ID, responseData.Categorys[0].ID)
		assert.Equal(t, expectedCategories[0].Name, responseData.Categorys[0].Name)
		assert.Equal(t, expectedCategories[1].ID, responseData.Categorys[1].ID)
		assert.Equal(t, expectedCategories[1].Name, responseData.Categorys[1].Name)
	})

	t.Run("success empty list", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetAllCategories(gomock.Any()).
			Return([]*models.Category{}, nil)

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		service.GetAllCategories(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseData dto.CategoryResponse
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		assert.NoError(t, err)

		assert.Equal(t, 0, responseData.Total)
		assert.Empty(t, responseData.Categorys)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetAllCategories(gomock.Any()).
			Return(nil, errors.New("database error"))

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		service.GetAllCategories(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestConvertToCategoriesResponse(t *testing.T) {
	t.Run("with categories", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()
		categories := []*models.Category{
			{
				ID:   id1,
				Name: "Category 1",
			},
			{
				ID:   id2,
				Name: "Category 2",
			},
		}

		result := dto.ConvertToCategoriesResponse(categories)

		assert.Equal(t, 2, result.Total)
		assert.Equal(t, categories[0].ID, result.Categorys[0].ID)
		assert.Equal(t, categories[0].Name, result.Categorys[0].Name)
		assert.Equal(t, categories[1].ID, result.Categorys[1].ID)
		assert.Equal(t, categories[1].Name, result.Categorys[1].Name)
	})

	t.Run("with nil category", func(t *testing.T) {
		id1 := uuid.New()
		categories := []*models.Category{
			nil,
			{
				ID:   id1,
				Name: "Category 1",
			},
			nil,
		}

		result := dto.ConvertToCategoriesResponse(categories)

		assert.Equal(t, 3, result.Total) // Total includes nil categories
		assert.Equal(t, 1, len(result.Categorys)) // Only non-nil categories are included
		assert.Equal(t, categories[1].ID, result.Categorys[0].ID)
		assert.Equal(t, categories[1].Name, result.Categorys[0].Name)
	})

	t.Run("empty list", func(t *testing.T) {
		categories := []*models.Category{}

		result := dto.ConvertToCategoriesResponse(categories)

		assert.Equal(t, 0, result.Total)
		assert.Empty(t, result.Categorys)
	})
}