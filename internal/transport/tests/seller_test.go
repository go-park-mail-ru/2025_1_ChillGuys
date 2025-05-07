package tests

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConvertToBriefProduct(t *testing.T) {
	t.Run("with seller info", func(t *testing.T) {
		productID := uuid.New()
		sellerID := uuid.New()

		product := &models.Product{
			ID:              productID,
			SellerID:        sellerID,
			Name:            "Test Product",
			Description:     "Test Description",
			Price:           100.0,
			PriceDiscount:   90.0,
			Quantity:        10,
			PreviewImageURL: "http://example.com/image.jpg",
			ReviewsCount:    5,
			Rating:          4.5,
			Seller: &models.Seller{
				Title:       "Test Seller",
				Description: "Test Seller Description",
			},
		}

		briefProduct := dto.ConvertToBriefProduct(product)

		assert.Equal(t, productID, briefProduct.ID)
		assert.Equal(t, "Test Product", briefProduct.Name)
		assert.Equal(t, "http://example.com/image.jpg", briefProduct.ImageURL)
		assert.Equal(t, 100.0, briefProduct.Price)
		assert.Equal(t, 90.0, briefProduct.PriceDiscount)
		assert.Equal(t, uint(10), briefProduct.Quantity)
		assert.Equal(t, uint(5), briefProduct.ReviewsCount)
		assert.Equal(t, float32(4.5), briefProduct.Rating)
		assert.NotNil(t, briefProduct.SellerInfo)
		assert.Equal(t, "Test Seller", briefProduct.SellerInfo.Title)
		assert.Equal(t, "Test Seller Description", briefProduct.SellerInfo.Description)
	})

	t.Run("without seller info", func(t *testing.T) {
		productID := uuid.New()
		sellerID := uuid.New()

		product := &models.Product{
			ID:              productID,
			SellerID:        sellerID,
			Name:            "Test Product",
			Description:     "Test Description",
			Price:           100.0,
			PriceDiscount:   90.0,
			Quantity:        10,
			PreviewImageURL: "http://example.com/image.jpg",
			ReviewsCount:    5,
			Rating:          4.5,
		}

		briefProduct := dto.ConvertToBriefProduct(product)

		assert.Equal(t, productID, briefProduct.ID)
		assert.Equal(t, "Test Product", briefProduct.Name)
		assert.Nil(t, briefProduct.SellerInfo)
	})
}

func TestConvertToProductsResponse(t *testing.T) {
	product1 := &models.Product{
		ID:              uuid.New(),
		Name:            "Product 1",
		PreviewImageURL: "http://example.com/image1.jpg",
		Price:           100.0,
		PriceDiscount:   90.0,
		Quantity:        10,
		ReviewsCount:    5,
		Rating:          4.5,
	}

	product2 := &models.Product{
		ID:              uuid.New(),
		Name:            "Product 2",
		PreviewImageURL: "http://example.com/image2.jpg",
		Price:           200.0,
		PriceDiscount:   180.0,
		Quantity:        20,
		ReviewsCount:    10,
		Rating:          4.0,
	}

	products := []*models.Product{product1, product2}

	response := dto.ConvertToProductsResponse(products)

	assert.Equal(t, 2, response.Total)
	assert.Equal(t, 2, len(response.Products))
	assert.Equal(t, "Product 1", response.Products[0].Name)
	assert.Equal(t, "Product 2", response.Products[1].Name)
}

func TestConvertToSellerProductsResponse(t *testing.T) {
	product1 := &models.Product{
		ID:              uuid.New(),
		Name:            "Product 1",
		PreviewImageURL: "http://example.com/image1.jpg",
		Price:           100.0,
		PriceDiscount:   90.0,
		Quantity:        10,
		ReviewsCount:    5,
		Rating:          4.5,
	}

	product2 := &models.Product{
		ID:              uuid.New(),
		Name:            "Product 2",
		PreviewImageURL: "http://example.com/image2.jpg",
		Price:           200.0,
		PriceDiscount:   180.0,
		Quantity:        20,
		ReviewsCount:    10,
		Rating:          4.0,
	}

	products := []*models.Product{product1, product2}

	response := dto.ConvertToSellerProductsResponse(products)

	assert.Equal(t, 2, response.Total)
	assert.Equal(t, 2, len(response.Products))
	assert.Equal(t, "Product 1", response.Products[0].Name)
	assert.Equal(t, "Product 2", response.Products[1].Name)
}

func TestDTOJsonMarshaling(t *testing.T) {
	t.Run("BriefProduct marshaling", func(t *testing.T) {
		product := dto.BriefProduct{
			ID:            uuid.New(),
			Name:          "Test Product",
			ImageURL:      "http://example.com/image.jpg",
			Price:         100.0,
			PriceDiscount: 90.0,
			Quantity:      10,
			ReviewsCount:  5,
			Rating:        4.5,
			SellerInfo: &dto.SellerInfo{
				Title:       "Test Seller",
				Description: "Test Description",
			},
		}

		data, err := json.Marshal(product)
		assert.NoError(t, err)

		var decoded dto.BriefProduct
		err = json.Unmarshal(data, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, product, decoded)
	})

	t.Run("ProductsResponse marshaling", func(t *testing.T) {
		response := dto.ProductsResponse{
			Total: 2,
			Products: []dto.BriefProduct{
				{
					ID:   uuid.New(),
					Name: "Product 1",
				},
				{
					ID:   uuid.New(),
					Name: "Product 2",
				},
			},
		}

		data, err := json.Marshal(response)
		assert.NoError(t, err)

		var decoded dto.ProductsResponse
		err = json.Unmarshal(data, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, response.Total, decoded.Total)
		assert.Equal(t, len(response.Products), len(decoded.Products))
	})
}
