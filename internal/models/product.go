package models

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	MediaFolder = "./media"
	CoverName = "cover.jpeg"
	DefaultPathCoverName = "product-default"
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

func GetProductCoverPath(id int) string {
	coverPath := filepath.Join(MediaFolder, fmt.Sprintf("product-%d", id), CoverName)

	if _, err := os.Stat(coverPath); os.IsNotExist(err) {
		// Если файл не существует, используем дефолтный URL
		coverPath = filepath.Join(MediaFolder, DefaultPathCoverName, CoverName)
	}

	return coverPath
}

func ConvertToBriefProduct(product *Product) BriefProduct{
	coverPath := GetProductCoverPath(product.ID)

	return BriefProduct{
		ID:           product.ID,
		Name:         product.Name,
		ImageURL:     coverPath,
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