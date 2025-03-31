package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
)

var testProducts = []*models.Product{
	{
		ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Name:        "Product 1",
		Description: "Description 1",
		Price:       1000,
		Quantity:    10,
		Rating:      4,
		ReviewsCount: 20,
	},
	{
		ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
		Name:        "Product 2",
		Description: "Description 2",
		Price:       2000,
		Quantity:    5,
		Rating:      4,
		ReviewsCount: 15,
	},
}

func TestGetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIProductUsecase(ctrl)
	logger := logrus.New()
	handler := transport.NewProductHandler(mockUsecase, logger)

	t.Run("Success case", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/products/", nil)
		assert.NoError(t, err)

		mockUsecase.EXPECT().
			GetAllProducts(req.Context()).
			Return(testProducts, nil).
			Times(1)

		rr := httptest.NewRecorder()
		handler.GetAllProducts(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response models.ProductsResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Products, len(testProducts))
		assert.Equal(t, len(testProducts), response.Total)
	})

	t.Run("Repository error case", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/products/", nil)
		assert.NoError(t, err)

		mockUsecase.EXPECT().
			GetAllProducts(req.Context()).
			Return(nil, errors.New("repository error")).
			Times(1)

		rr := httptest.NewRecorder()
		handler.GetAllProducts(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Failed get all products")
	})
}

func TestGetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIProductUsecase(ctrl)
	logger := logrus.New()
	handler := transport.NewProductHandler(mockUsecase, logger)

	t.Run("Success case", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		testProduct := &models.Product{
			ID:          testID,
			Name:        "Смартфон Xiaomi Redmi Note 10",
			Description: "Смартфон с AMOLED-дисплеем и камерой 48 Мп",
			Quantity:    50,
			Price:       19999,
			ReviewsCount: 120,
			Rating:      4,
		}

		req, err := http.NewRequest("GET", "/products/"+testID.String(), nil)
		assert.NoError(t, err)

		vars := map[string]string{
			"id": testID.String(),
		}
		req = mux.SetURLVars(req, vars)

		mockUsecase.EXPECT().
			GetProductByID(req.Context(), testID).
			Return(testProduct, nil).
			Times(1)

		rr := httptest.NewRecorder()
		handler.GetProductByID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var product models.Product
		err = json.Unmarshal(rr.Body.Bytes(), &product)
		assert.NoError(t, err)
		assert.Equal(t, testProduct.ID, product.ID)
		assert.Equal(t, testProduct.Name, product.Name)
	})

	t.Run("Invalid ID case", func(t *testing.T) {
		invalidID := "abc"

		req, err := http.NewRequest("GET", "/products/"+invalidID, nil)
		assert.NoError(t, err)

		vars := map[string]string{
			"id": invalidID,
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler.GetProductByID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid ID")
	})

	t.Run("Not found case", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440999")

		req, err := http.NewRequest("GET", "/products/"+testID.String(), nil)
		assert.NoError(t, err)

		vars := map[string]string{
			"id": testID.String(),
		}
		req = mux.SetURLVars(req, vars)

		mockUsecase.EXPECT().
			GetProductByID(req.Context(), testID).
			Return(nil, errors.New("not found")).
			Times(1)

		rr := httptest.NewRecorder()
		handler.GetProductByID(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestGetProductCover(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockIProductUsecase(ctrl)
	logger := logrus.New()
	handler := transport.NewProductHandler(mockUsecase, logger)

	t.Run("Success case", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		testCoverData := []byte{0xFF, 0xD8, 0xFF}

		req, err := http.NewRequest("GET", "/products/"+testID.String()+"/cover", nil)
		assert.NoError(t, err)

		vars := map[string]string{
			"id": testID.String(),
		}
		req = mux.SetURLVars(req, vars)

		mockUsecase.EXPECT().
			GetProductCover(req.Context(), testID).
			Return(testCoverData, nil).
			Times(1)

		rr := httptest.NewRecorder()
		handler.GetProductCover(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "image/jpeg", rr.Header().Get("Content-Type"))
		assert.Equal(t, testCoverData, rr.Body.Bytes())
	})

	t.Run("Cover not found case", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

		req, err := http.NewRequest("GET", "/products/"+testID.String()+"/cover", nil)
		assert.NoError(t, err)

		vars := map[string]string{
			"id": testID.String(),
		}
		req = mux.SetURLVars(req, vars)

		mockUsecase.EXPECT().
			GetProductCover(req.Context(), testID).
			Return(nil, os.ErrNotExist).
			Times(1)

		rr := httptest.NewRecorder()
		handler.GetProductCover(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Cover file not found")
	})

	t.Run("Invalid ID case", func(t *testing.T) {
		invalidID := "abc"

		req, err := http.NewRequest("GET", "/products/"+invalidID+"/cover", nil)
		assert.NoError(t, err)

		vars := map[string]string{
			"id": invalidID,
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler.GetProductCover(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid ID")
	})

	t.Run("Internal server error case", func(t *testing.T) {
		testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")

		req, err := http.NewRequest("GET", "/products/"+testID.String()+"/cover", nil)
		assert.NoError(t, err)

		vars := map[string]string{
			"id": testID.String(),
		}
		req = mux.SetURLVars(req, vars)

		mockUsecase.EXPECT().
			GetProductCover(req.Context(), testID).
			Return(nil, errors.New("internal error")).
			Times(1)

		rr := httptest.NewRecorder()
		handler.GetProductCover(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Failed to get cover file")
	})
}