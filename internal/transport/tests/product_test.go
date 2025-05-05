package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	minio_mocks "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func createTestContext() context.Context {
	return logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
}

func TestProductService_GetAllProducts(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUsecase := mocks.NewMockIProductUsecase(ctrl)
    mockMinio := minio_mocks.NewMockProvider(ctrl)

    service := product.NewProductService(mockUsecase, mockMinio)

    testProducts := []*models.Product{
        {
            ID:              uuid.New(),
            Name:            "Product 1",
            PreviewImageURL: "url1",
            Price:           10.5,
            Status:          models.ProductApproved,
        },
    }

    t.Run("Success with default offset", func(t *testing.T) {
        mockUsecase.EXPECT().
            GetAllProducts(gomock.Any(), 0).
            Return(testProducts, nil)

        req := httptest.NewRequest("GET", "/products", nil)
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        // Set up router with route that has optional offset parameter
        router := mux.NewRouter()
        router.HandleFunc("/products", service.GetAllProducts).Methods("GET")
        router.HandleFunc("/products/{offset}", service.GetAllProducts).Methods("GET")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var response dto.ProductsResponse
        err := json.NewDecoder(resp.Body).Decode(&response)
        assert.NoError(t, err)
        assert.Equal(t, len(testProducts), len(response.Products))
    })

    t.Run("Success with offset", func(t *testing.T) {
        offset := 10
        mockUsecase.EXPECT().
            GetAllProducts(gomock.Any(), offset).
            Return(testProducts, nil)

        req := httptest.NewRequest("GET", "/products/"+strconv.Itoa(offset), nil)
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        // Set up router with route that has optional offset parameter
        router := mux.NewRouter()
        router.HandleFunc("/products", service.GetAllProducts).Methods("GET")
        router.HandleFunc("/products/{offset}", service.GetAllProducts).Methods("GET")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var response dto.ProductsResponse
        err := json.NewDecoder(resp.Body).Decode(&response)
        assert.NoError(t, err)
        assert.Equal(t, len(testProducts), len(response.Products))
    })
}

func TestProductService_GetProductByID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUsecase := mocks.NewMockIProductUsecase(ctrl)
    mockMinio := minio_mocks.NewMockProvider(ctrl)

    service := product.NewProductService(mockUsecase, mockMinio)

    testID := uuid.New()

    t.Run("Success", func(t *testing.T) {
        mockUsecase.EXPECT().
            GetProductByID(gomock.Any(), testID).
            Return(&models.Product{
                ID:              testID,
                Name:            "Test Product",
                PreviewImageURL: "test_url",
                Price:           15.99,
                Status:          models.ProductApproved,
            }, nil)

        req := httptest.NewRequest("GET", "/products/"+testID.String(), nil)
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/{id}", service.GetProductByID).Methods("GET")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusOK, resp.StatusCode)
    })

    t.Run("Invalid ID", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/products/invalid", nil)
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/{id}", service.GetProductByID).Methods("GET")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })
}

func TestProductService_CreateOne(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUsecase := mocks.NewMockIProductUsecase(ctrl)
    mockMinio := minio_mocks.NewMockProvider(ctrl)

    service := product.NewProductService(mockUsecase, mockMinio)

    t.Run("Success", func(t *testing.T) {
        mockMinio.EXPECT().
            CreateOne(gomock.Any(), gomock.Any()).
            Return(&dto.UploadResponse{URL: "http://test.url"}, nil)

        body := &bytes.Buffer{}
        writer := multipart.NewWriter(body)
        part, _ := writer.CreateFormFile("file", "test.jpg")
        part.Write([]byte("test image content"))
        writer.Close()

        req := httptest.NewRequest("POST", "/products/upload", body)
        req.Header.Set("Content-Type", writer.FormDataContentType())
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/upload", service.CreateOne).Methods("POST")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusOK, resp.StatusCode)
    })
}

func TestProductService_GetProductsByCategory(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUsecase := mocks.NewMockIProductUsecase(ctrl)
    mockMinio := minio_mocks.NewMockProvider(ctrl)

    service := product.NewProductService(mockUsecase, mockMinio)

    categoryID := uuid.New()
    testProducts := []*models.Product{
        {
            ID:              uuid.New(),
            Name:            "Test Product",
            PreviewImageURL: "test_url",
            Price:           15.99,
            Status:          models.ProductApproved,
        },
    }

    t.Run("Success with default params", func(t *testing.T) {
        mockUsecase.EXPECT().
            GetProductsByCategory(
                gomock.Any(), 
                categoryID, 
                0, 
                0.0,  // minPrice as float64
                0.0,  // maxPrice as float64
                float32(0),  // minRating as float32
                models.SortByDefault,
            ).
            Return(testProducts, nil)

        req := httptest.NewRequest("GET", "/products/category/"+categoryID.String(), nil)
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/category/{id}", service.GetProductsByCategory).Methods("GET")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusOK, resp.StatusCode)
    })

    t.Run("Invalid category ID", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/products/category/invalid", nil)
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/category/{id}", service.GetProductsByCategory).Methods("GET")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })
}

func TestProductService_AddProduct(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUsecase := mocks.NewMockIProductUsecase(ctrl)
    mockMinio := minio_mocks.NewMockProvider(ctrl)
    service := product.NewProductService(mockUsecase, mockMinio)

    // Test data
    validCategoryID := uuid.New()
    validSellerID := uuid.New()
    createdProductID := uuid.New()

    t.Run("Success", func(t *testing.T) {
		// Setup mock response with proper ProductStatus
		mockResponse := &models.Product{
			ID:              createdProductID,
			SellerID:        validSellerID,
			Name:           "Test Product",
			PreviewImageURL: "http://test.com/image.jpg",
			Description:    "Test description",
			Price:          19.99,
			PriceDiscount:   15.99,
			Quantity:       100,
			Rating:         4.5,
			ReviewsCount:   10,
			Status:         models.ProductApproved, // Make sure this matches your actual ProductStatus type
		}
	
		mockUsecase.EXPECT().
			AddProduct(gomock.Any(), gomock.Any(), validCategoryID).
			Return(mockResponse, nil)
	
		// Create request body without status field since it's not in AddProductRequest
		requestBody := dto.AddProductRequest{
			Name:            "Test Product",
			SellerID:        validSellerID.String(),
			PreviewImageURL: "http://test.com/image.jpg",
			Description:     "Test description",
			Price:           19.99,
			PriceDiscount:   15.99,
			Quantity:        100,
			Rating:         4.5,
			ReviewsCount:    10,
			Category:       validCategoryID.String(),
		}
	
		body, err := json.Marshal(requestBody)
		assert.NoError(t, err)
	
		req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(createTestContext())
		w := httptest.NewRecorder()
	
		router := mux.NewRouter()
		router.HandleFunc("/products", service.AddProduct).Methods("POST")
		router.ServeHTTP(w, req)
	
		resp := w.Result()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	
		// Create a proper response struct that matches what your handler actually returns
		var response struct {
			ID              uuid.UUID `json:"id"`
			Name            string    `json:"name"`
			PreviewImageURL string    `json:"preview_image_url"`
			Description     string    `json:"description"`
			Price           float64   `json:"price"`
			PriceDiscount   float64   `json:"price_discount"`
			Quantity        uint      `json:"quantity"`
			Rating          float32   `json:"rating"`
			ReviewsCount    uint      `json:"reviews_count"`
		}
	
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, createdProductID, response.ID)
		assert.Equal(t, "Test Product", response.Name)
	})

    t.Run("Invalid Category ID", func(t *testing.T) {
        requestBody := map[string]interface{}{
            "name":      "Test Product",
            "seller_id": validSellerID.String(),
            "category":  "invalid-uuid",
            "price":     19.99,
            "quantity":  100,
        }

        body, _ := json.Marshal(requestBody)
        req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products", service.AddProduct).Methods("POST")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("Validation Error - Missing Name", func(t *testing.T) {
        mockUsecase.EXPECT().
            AddProduct(gomock.Any(), gomock.Any(), validCategoryID).
            Return(nil, errs.ErrEmptyProductName)

        requestBody := map[string]interface{}{
            "seller_id": validSellerID.String(),
            "category":  validCategoryID.String(),
            "price":     19.99,
            "quantity":  100,
        }

        body, _ := json.Marshal(requestBody)
        req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products", service.AddProduct).Methods("POST")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })
}

func TestProductService_GetProductsByIDs(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUsecase := mocks.NewMockIProductUsecase(ctrl)
    mockMinio := minio_mocks.NewMockProvider(ctrl)
    service := product.NewProductService(mockUsecase, mockMinio)

    // Test data
    productID1 := uuid.New()
    productID2 := uuid.New()
    testProducts := []*models.Product{
        {
            ID:              productID1,
            Name:            "Product 1",
            PreviewImageURL: "url1",
            Price:           10.5,
            Status:          models.ProductApproved,
        },
        {
            ID:              productID2,
            Name:            "Product 2",
            PreviewImageURL: "url2",
            Price:           20.5,
            Status:          models.ProductApproved,
        },
    }

    t.Run("Success", func(t *testing.T) {
        mockUsecase.EXPECT().
            GetProductsByIDs(gomock.Any(), []uuid.UUID{productID1, productID2}).
            Return(testProducts, nil)

        requestBody := dto.GetProductsByIDRequest{
            ProductIDs: []uuid.UUID{productID1, productID2},
        }

        body, _ := json.Marshal(requestBody)
        req := httptest.NewRequest("POST", "/products/batch", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/batch", service.GetProductsByIDs).Methods("POST")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var response dto.ProductsResponse
        err := json.NewDecoder(resp.Body).Decode(&response)
        assert.NoError(t, err)
        assert.Equal(t, 2, len(response.Products))
        assert.Equal(t, productID1, response.Products[0].ID)
        assert.Equal(t, productID2, response.Products[1].ID)
    })

    t.Run("Empty request body", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/products/batch", nil)
        req.Header.Set("Content-Type", "application/json")
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/batch", service.GetProductsByIDs).Methods("POST")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("Invalid JSON", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/products/batch", bytes.NewReader([]byte("invalid json")))
        req.Header.Set("Content-Type", "application/json")
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/batch", service.GetProductsByIDs).Methods("POST")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("Empty product IDs", func(t *testing.T) {
        requestBody := dto.GetProductsByIDRequest{
            ProductIDs: []uuid.UUID{},
        }

        body, _ := json.Marshal(requestBody)
        req := httptest.NewRequest("POST", "/products/batch", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        req = req.WithContext(createTestContext())
        w := httptest.NewRecorder()

        router := mux.NewRouter()
        router.HandleFunc("/products/batch", service.GetProductsByIDs).Methods("POST")
        router.ServeHTTP(w, req)

        resp := w.Result()
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

        var errorResponse dto.ErrorResponseDTO
        err := json.NewDecoder(resp.Body).Decode(&errorResponse)
        assert.NoError(t, err)
        assert.Contains(t, errorResponse.Message, "at least one product ID is required")
    })
}