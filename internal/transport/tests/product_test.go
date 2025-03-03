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

func TestGetProductByID(t *testing.T){
	
}