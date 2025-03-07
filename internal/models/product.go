package models

import (
	"fmt"
	"os"
	"path/filepath"
)

type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Count        uint    `json:"count"`
	Price        uint    `json:"price"`
	ReviewsCount uint    `json:"reviews_count"`
	Rating       float64 `json:"rating"`
}

type BriefProduct struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	ImageURL     string  `json:"image"`
	Price        uint    `json:"price"`
	ReviewsCount uint    `json:"reviews_count"`
	Rating       float64 `json:"rating"`
}

func ConvertToBriefProduct(product *Product) BriefProduct{
	coverPath := filepath.Join("./media", fmt.Sprintf("product-%s", product.ID), "cover.jpeg")

	imageURL := fmt.Sprintf("media/product-%s/cover.jpeg", product.ID)

	if _, err := os.Stat(coverPath); os.IsNotExist(err) {
		// Если файл не существует, используем дефолтный URL
		imageURL = "media/product-default/cover.jpeg"
	}

	return BriefProduct{
		ID:           product.ID,
		Name:         product.Name,
		ImageURL:     imageURL,
		Price:        product.Price,
		ReviewsCount: product.ReviewsCount,
		Rating:       product.Rating,
	}
}

type ProductsResponse struct {
	Total    int                   `json:"total"`
	Products []BriefProduct 		   `json:"products"`
}

func ConvertToProductsResponse(products []*Product) ProductsResponse {
	briefProducts := make([]BriefProduct, 0, len(products))
	for _, product := range products {
		briefProduct := ConvertToBriefProduct(product)
		briefProducts = append(briefProducts, briefProduct)
	}

	response := ProductsResponse{
		Total: len(briefProducts),
		Products: briefProducts,
	}

	return response
}