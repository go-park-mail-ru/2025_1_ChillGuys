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
		ImageURL:     fmt.Sprintf("media/products/product-%d.jpg", product.ID),
		Price:        product.Price,
		ReviewsCount: product.ReviewsCount,
		Rating:       product.Rating,
	}
}