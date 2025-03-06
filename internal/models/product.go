package models

import "fmt"

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
	return BriefProduct{
		ID:           product.ID,
		Name:         product.Name,
		ImageURL:     fmt.Sprintf("media/products/product-%d.jpeg", product.ID),
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