package tests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	minio_mocks "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// ProductResponse - структура для парсинга ответа с продуктом
type ProductResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	PreviewImageURL string    `json:"preview_image_url"`
	Price           float64   `json:"price"`
	ReviewsCount    uint      `json:"reviews_count"`
	Rating          uint      `json:"rating"`
}

func parseResponseBody(t *testing.T, body io.ReadCloser, target interface{}) {
	t.Helper()
	err := json.NewDecoder(body).Decode(target)
	assert.NoError(t, err)
}

func TestProductService_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIProductUsecase(ctrl)
	mockMinio := minio_mocks.NewMockProvider(ctrl)

	service := product.NewProductService(mockUsecase, mockMinio)

	t.Run("success", func(t *testing.T) {
		expectedProducts := []*models.Product{
			{
				ID:              uuid.New(),
				Name:            "Product 1",
				PreviewImageURL: "url1",
				Price:           10.5,
				ReviewsCount:    5,
				Rating:          4,
				Status:          models.ProductApproved,
			},
			{
				ID:              uuid.New(),
				Name:            "Product 2",
				PreviewImageURL: "url2",
				Price:           20.5,
				ReviewsCount:    10,
				Rating:          5,
				Status:          models.ProductApproved,
			},
		}

		mockUsecase.EXPECT().
			GetAllProducts(gomock.Any()).
			Return(expectedProducts, nil)

		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()

		service.GetAllProducts(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseData dto.ProductsResponse
		parseResponseBody(t, resp.Body, &responseData)

		assert.Equal(t, len(expectedProducts), responseData.Total)
		assert.Equal(t, expectedProducts[0].Name, responseData.Products[0].Name)
		assert.Equal(t, expectedProducts[0].PreviewImageURL, responseData.Products[0].ImageURL)
		assert.Equal(t, expectedProducts[0].Price, responseData.Products[0].Price)
		assert.Equal(t, expectedProducts[0].ReviewsCount, responseData.Products[0].ReviewsCount)
		assert.Equal(t, expectedProducts[0].Rating, responseData.Products[0].Rating)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetAllProducts(gomock.Any()).
			Return(nil, errors.New("some error"))

		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()

		service.GetAllProducts(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestProductService_GetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIProductUsecase(ctrl)
	mockMinio := minio_mocks.NewMockProvider(ctrl)

	service := product.NewProductService(mockUsecase, mockMinio)

	t.Run("success", func(t *testing.T) {
		productID := uuid.New()
		expectedProduct := &models.Product{
			ID:              productID,
			Name:            "Test Product",
			PreviewImageURL: "test_url",
			Price:           15.99,
			ReviewsCount:    8,
			Rating:          4,
			Status:          models.ProductApproved,
		}

		mockUsecase.EXPECT().
			GetProductByID(gomock.Any(), productID).
			Return(expectedProduct, nil)

		req := httptest.NewRequest("GET", "/products/"+productID.String(), nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.GetProductByID(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseData ProductResponse
		parseResponseBody(t, resp.Body, &responseData)

		assert.Equal(t, expectedProduct.ID, responseData.ID)
		assert.Equal(t, expectedProduct.Name, responseData.Name)
		assert.Equal(t, expectedProduct.PreviewImageURL, responseData.PreviewImageURL)
		assert.Equal(t, expectedProduct.Price, responseData.Price)
		assert.Equal(t, expectedProduct.ReviewsCount, responseData.ReviewsCount)
		assert.Equal(t, expectedProduct.Rating, responseData.Rating)
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/products/invalid", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

		service.GetProductByID(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		productID := uuid.New()

		mockUsecase.EXPECT().
			GetProductByID(gomock.Any(), productID).
			Return(nil, errs.NewNotFoundError("product not found"))

		req := httptest.NewRequest("GET", "/products/"+productID.String(), nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

		service.GetProductByID(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestProductService_GetProductsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIProductUsecase(ctrl)
	mockMinio := minio_mocks.NewMockProvider(ctrl)

	service := product.NewProductService(mockUsecase, mockMinio)

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/products/category/invalid", nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

		service.GetProductsByCategory(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		categoryID := uuid.New()

		mockUsecase.EXPECT().
			GetProductsByCategory(gomock.Any(), categoryID).
			Return(nil, errs.NewNotFoundError("category not found"))

		req := httptest.NewRequest("GET", "/api/v1/products/category/"+categoryID.String(), nil)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": categoryID.String()})

		service.GetProductsByCategory(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
