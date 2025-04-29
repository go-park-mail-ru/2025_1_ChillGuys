package helpers

import "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"

func GetFinalPrice(p *models.Product) float64 {
    if p.PriceDiscount > 0 {
        return p.PriceDiscount
    }
    return p.Price
}