package models

type CategorySuggestion struct {
	Name string `json:"name" db:"name"`
}

type ProductSuggestion struct {
	Name string `json:"name" db:"name"`
}
