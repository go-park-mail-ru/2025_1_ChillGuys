package dto

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
)

type BriefProduct struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ImageURL      string    `json:"image"`
	Price         float64   `json:"price"`
	PriceDiscount float64   `json:"discount_price"`
	Quantity      uint      `json:"quantity"`
	ReviewsCount  uint      `json:"reviews_count"`
	Rating        float32   `json:"rating"`
}

func ConvertToBriefProduct(product *models.Product) BriefProduct {
	return BriefProduct{
		ID:            product.ID,
		Name:          product.Name,
		ImageURL:      product.PreviewImageURL,
		Price:         product.Price,
		PriceDiscount: product.PriceDiscount,
		Quantity:      product.Quantity,
		ReviewsCount:  product.ReviewsCount,
		Rating:        product.Rating,
	}
}

type ProductsResponse struct {
	Total    int            `json:"total"`
	Products []BriefProduct `json:"products"`
}

func ConvertToProductsResponse(products []*models.Product) ProductsResponse {
	briefProducts := make([]BriefProduct, 0, len(products))
	for _, product := range products {
		briefProducts = append(briefProducts, ConvertToBriefProduct(product))
	}

	return ProductsResponse{
		Total:    len(briefProducts),
		Products: briefProducts,
	}
}

type GetProductsByIDRequest struct {
	ProductIDs []uuid.UUID `json:"productIDs" validate:"required,min=1"`
}
