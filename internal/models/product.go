package models

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const (
	MediaFolder          = "./media"
	CoverName            = "cover.jpeg"
	DefaultPathCoverName = "product-default"
)

type Product struct {
	ID              uuid.UUID `json:"id" db:"id"`
	SellerID        uuid.UUID `json:"seller_id" db:"seller_id"`
	Name            string    `json:"name" db:"name"`
	PreviewImageURL string    `json:"preview_image_url,omitempty" db:"preview_image_url"`
	Description     string    `json:"description,omitempty" db:"description"`
	Status          string    `json:"status" db:"status"`
	Price           uint      `json:"price" db:"price"`
	Quantity        uint      `json:"quantity" db:"quantity"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	Rating          int       `json:"rating,omitempty" db:"rating"`
	ReviewsCount    uint      `json:"reviews_count" db:"reviews_count"` // Добавлено поле
}

type BriefProduct struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ImageURL     string    `json:"image"`
	Price        uint      `json:"price"`
	ReviewsCount uint      `json:"reviews_count"`
	Rating       int       `json:"rating"`
}

func GetProductCoverPath(id uuid.UUID) string {
	coverPath := filepath.Join(MediaFolder, fmt.Sprintf("product-%d", id), CoverName)

	if _, err := os.Stat(coverPath); os.IsNotExist(err) {
		// Если файл не существует, используем дефолтный URL
		coverPath = "http://minio:9000/bazaar-bucket/84cf2b13-318b-48e3-88f6-30abf86d4a6b?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20250401%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20250401T193024Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=c3a6ef3254747741cc27459b428f981d15e3cc4ec92914efae6f7d65dc911c85"
	}

	return coverPath
}

func ConvertToBriefProduct(product *Product) BriefProduct {
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
	Total    int            `json:"total"`
	Products []BriefProduct `json:"products"`
}

func ConvertToProductsResponse(products []*Product) ProductsResponse {
	briefProducts := make([]BriefProduct, 0, len(products))
	for _, product := range products {
		briefProducts = append(briefProducts, ConvertToBriefProduct(product))
	}

	return ProductsResponse{
		Total:    len(briefProducts),
		Products: briefProducts,
	}
}
