package models

import (
	"time"

	"github.com/google/uuid"
)

type BasketItem struct {
	ID             uuid.UUID  `json:"id"`
	BasketID       uuid.UUID  `json:"basket_id"`
	ProductID      uuid.UUID  `json:"product_id"`
	Quantity       int        `json:"quantity"`
	UpdatedAt      time.Time  `json:"updated_at"`
	ProductName    string     `json:"product_name"`
	Price		   float64    `json:"product_price"`
	ProductImage   string     `json:"product_image"`
	PriceDiscount  float64    `json:"price_discount"`
}

type Basket struct {
	ID                  uuid.UUID  `json:"id"`
	UserID              uuid.UUID  `json:"user_id"`
	TotalPrice          float64    `json:"total_price"`
	TotalPriceDiscount  float64    `json:"total_price_discount"`
}

type BasketResponse struct {
    Total              int          `json:"total"`
    TotalPrice         float64      `json:"total_price"`
    TotalPriceDiscount float64      `json:"total_price_discount"`
	Products           []BasketItem `json:"products"`
}

func ConvertToBasketResponse(items []*BasketItem) BasketResponse{
	var totalPrice, totalPriceDiscount float64 
	productsList := make([]BasketItem, 0, len(items))
	for _, product := range items{
		productsList = append(productsList, *product)
		totalPrice = totalPrice + product.Price * float64(product.Quantity)
        price := product.Price
        if product.PriceDiscount > 0 {
            price = product.PriceDiscount
        }
        totalPriceDiscount += price * float64(product.Quantity)
	}

	return BasketResponse{
		Total: len(productsList),
		TotalPrice: totalPrice,
		TotalPriceDiscount : totalPriceDiscount,
		Products: productsList,
	}
}