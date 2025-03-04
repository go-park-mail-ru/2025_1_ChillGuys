package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProducts(t *testing.T) {
    // Инициализация репозитория и обработчика
    repo := repository.NewProductRepo()
    handler := transport.NewProductHandler(repo)

    // Создание HTTP-запроса
    req, err := http.NewRequest("GET", "/products", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    // Создание ResponseRecorder для записи ответа
    rr := httptest.NewRecorder()

    // Вызов обработчика
    handler.GetAllProducts(rr, req)

    // Проверка статус-кода
    assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200, got %d", rr.Code)

    // Декодирование JSON-ответа
    var response struct {
        Total    int             `json:"total"`
        Products []models.BriefProduct `json:"products"`
    }
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err, "Failed to decode response")

    // Проверка, что список продуктов не пуст
    assert.NotEmpty(t, response.Products, "Expected non-empty product list")
}

func TestGetProductByID(t *testing.T) {
    // Инициализация репозитория и обработчика
    repo := repository.NewProductRepo()
    handler := transport.NewProductHandler(repo)

    // Тестовый ID продукта
    testID := 1

    // Создание HTTP-запроса с path-параметром
    req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(testID), nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    // Добавление path-параметра в запрос
    vars := map[string]string{
        "id": strconv.Itoa(testID),
    }
    req = mux.SetURLVars(req, vars)

    // Создание ResponseRecorder для записи ответа
    rr := httptest.NewRecorder()

    // Вызов обработчика
    handler.GetProductByID(rr, req)

    // Проверка статус-кода
    assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200, got %d", rr.Code)

    // Декодирование JSON-ответа
    var product models.Product
    err = json.Unmarshal(rr.Body.Bytes(), &product)
    assert.NoError(t, err, "Failed to decode response")

    // Проверка данных продукта
    assert.Equal(t, testID, product.ID, "Expected product ID %d, got %d", testID, product.ID)
    assert.NotEmpty(t, product.Name, "Product name should not be empty")
    assert.NotEmpty(t, product.Description, "Product description should not be empty")
    assert.Greater(t, product.Price, uint(0), "Product price should be greater than 0")
    assert.GreaterOrEqual(t, product.Rating, 0.0, "Product rating should be greater than or equal to 0")
}