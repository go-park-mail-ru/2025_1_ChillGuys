package dto

import "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"

type SearchReq struct {
	SubString string `json:"sub_string"`
}

type SearchProductsByNamesReq struct {
	ProductNames []models.ProductSuggestion `json:"product_names"`
}

type SearchResponse struct {
	Categories CategoryResponse `json:"categories"`
	Products   ProductsResponse `json:"products"`
}
