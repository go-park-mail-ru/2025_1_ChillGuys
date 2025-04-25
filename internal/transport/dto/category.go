package dto

import "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"

type CategoryResponse struct {
	Total     int               `json:"total"`
	Categorys []models.Category `json:"categories"`
}

// ConvertToCategoriesResponse - преобразует список категорий в структуру CategoryResponse.
func ConvertToCategoriesResponse(categories []*models.Category) CategoryResponse {
	categoryList := make([]models.Category, 0, len(categories))
	for _, cat := range categories {
		if cat != nil {
			categoryList = append(categoryList, *cat)
		}
	}

	return CategoryResponse{
		Total:     len(categories),
		Categorys: categoryList,
	}
}
