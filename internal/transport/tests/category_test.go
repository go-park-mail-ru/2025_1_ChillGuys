package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/category"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	t.Run("success", func(t *testing.T) {
		categories := []*models.Category{
			{ID: uuid.New(), Name: "Электроника"},
			{ID: uuid.New(), Name: "Одежда"},
		}

		mockUsecase.EXPECT().
			GetAllCategories(gomock.Any()).
			Return(categories, nil)

		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
		w := httptest.NewRecorder()

		service.GetAllCategories(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result dto.CategoryResponse
		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, "Электроника", result.Categorys[0].Name)
	})

	t.Run("internal error", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetAllCategories(gomock.Any()).
			Return(nil, errors.New("some error"))

		req := httptest.NewRequest("GET", "/api/v1/categories", nil)
		w := httptest.NewRecorder()

		service.GetAllCategories(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestCategoryService_GetAllSubcategories(t *testing.T) {
	mockUsecase, service := setupTestCategory(t)
	categoryID := uuid.New()

	t.Run("success", func(t *testing.T) {
		subcategories := []*models.Category{
			{ID: uuid.New(), Name: "Смартфоны"},
			{ID: uuid.New(), Name: "Ноутбуки"},
		}

		mockUsecase.EXPECT().
			GetAllSubategories(gomock.Any(), categoryID).
			Return(subcategories, nil)

		req := httptest.NewRequest("GET", "/api/v1/categories/"+categoryID.String()+"/subcategories", nil)
		req = mux.SetURLVars(req, map[string]string{"id": categoryID.String()})
		w := httptest.NewRecorder()

		service.GetAllSubcategories(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result dto.CategoryResponse
		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, 2, result.Total)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/categories/invalid-id/subcategories", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})
		w := httptest.NewRecorder()

		service.GetAllSubcategories(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestCategoryService_GetNameSubcategory(t *testing.T) {
	mockUsecase, service := setupTestCategory(t)
	subcategoryID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetNameSubcategory(gomock.Any(), subcategoryID).
			Return("Планшеты", nil)

		req := httptest.NewRequest("GET", "/api/v1/subcategories/"+subcategoryID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"id": subcategoryID.String()})
		w := httptest.NewRecorder()

		service.GetNameSubcategory(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result dto.NameSubcategory
		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, "Планшеты", result.Name)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/subcategories/invalid-id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})
		w := httptest.NewRecorder()

		service.GetNameSubcategory(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("subcategory not found", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetNameSubcategory(gomock.Any(), subcategoryID).
			Return("", errs.ErrNotFound)

		req := httptest.NewRequest("GET", "/api/v1/subcategories/"+subcategoryID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"id": subcategoryID.String()})
		w := httptest.NewRecorder()

		service.GetNameSubcategory(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
