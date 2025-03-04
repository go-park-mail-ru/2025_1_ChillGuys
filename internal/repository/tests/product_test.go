package tests

import (
	"testing"
	"context"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProducts(t *testing.T){
	repo := repository.NewProductRepo()

	products, err := repo.GetAllProducts(context.Background())

	assert.NoError(t, err)

	expectedCount := 15 // Мы заполнили репозиторий 15 продуктами в populateMockData
    assert.Equal(t, expectedCount, len(products), "Expected %d products, got %d", expectedCount, len(products))
}

func TestGetProductByID(t *testing.T) {
    repo := repository.NewProductRepo()

    t.Run("Existing product", func(t *testing.T) {
        testID := 1

        product, err := repo.GetProductByID(context.Background(), testID)

        // Проверка, что ошибки нет
        assert.NoError(t, err, "GetProductByID should not return an error")

        // Проверка, что продукт не nil
        assert.NotNil(t, product, "Product should not be nil")

        // Проверка ID продукта
        assert.Equal(t, testID, product.ID, "Expected product ID %d, got %d", testID, product.ID)

        // Проверка данных продукта
        assert.NotEmpty(t, product.Name, "Product name should not be empty")
        assert.NotEmpty(t, product.Description, "Product description should not be empty")
        assert.Greater(t, product.Price, uint(0), "Product price should be greater than 0")
        assert.GreaterOrEqual(t, product.Rating, 0.0, "Product rating should be greater than or equal to 0")
    })

    // Тест для несуществующего продукта
    t.Run("Non-existing product", func(t *testing.T) {
        nonExistingID := 999

        product, err := repo.GetProductByID(context.Background(), nonExistingID)

        // Проверка, что ошибка есть
        assert.Error(t, err, "GetProductByID should return an error for non-existing product")

        // Проверка, что продукт nil
        assert.Nil(t, product, "Product should be nil for non-existing ID")

        // Проверка текста ошибки
        assert.Contains(t, err.Error(), "not found", "Error message should contain 'not found'")
    })
}