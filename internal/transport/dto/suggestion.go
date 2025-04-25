package dto

import "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"

type CategoryNameResponse struct {
	CategoriesNames []models.CategorySuggestion `json:"categories_names"`
}

type ProductNameResponse struct {
	ProductNames []models.ProductSuggestion `json:"product_names"`
}

type SuggestionsReq struct {
	SubString string `json:"sub_string"`
}

type CombinedSuggestionsResponse struct {
	Categories []models.CategorySuggestion `json:"categories"`
	Products   []models.ProductSuggestion  `json:"products"`
}
