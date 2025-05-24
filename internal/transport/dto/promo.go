package dto

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
)

type CreatePromoRequest struct {
	Code      string    `json:"code" validate:"required"`
	Percent   int       `json:"percent" validate:"required,min=1,max=100"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type PromoResponse struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Percent   int       `json:"percent"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type PromosResponse struct {
	Total int             `json:"total"`
	Promos []PromoResponse `json:"promos"`
}

func ConvertToPromoResponse(promo *models.PromoCode) PromoResponse {
	return PromoResponse{
		ID:        promo.ID,
		Code:      promo.Code,
		Percent:   promo.Percent,
		StartDate: promo.StartDate,
		EndDate:   promo.EndDate,
	}
}

func ConvertToPromosResponse(promos []*models.PromoCode) PromosResponse {
	promoResponses := make([]PromoResponse, 0, len(promos))
	for _, promo := range promos {
		promoResponses = append(promoResponses, ConvertToPromoResponse(promo))
	}

	return PromosResponse{
		Total:  len(promoResponses),
		Promos: promoResponses,
	}
}

type CheckPromoRequest struct {
    Code string `json:"code" validate:"required"`
}

type PromoValidityResponse struct {
    IsValid bool `json:"is_valid"`
    Percent int  `json:"percent,omitempty"` 
}