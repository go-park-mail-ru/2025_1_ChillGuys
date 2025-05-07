package dto

import (
	"github.com/guregu/null"
)

type SearchReq struct {
	CategoryID null.String `json:"category_id"`
	SubString  string      `json:"sub_string"`
}

type SearchResponse struct {
	Categories CategoryResponse `json:"categories"`
	Products   ProductsResponse `json:"products"`
}
