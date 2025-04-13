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
	QuantityRemain int  	  `json:"remain_quantity"`
}

type Basket struct {
	ID                  uuid.UUID  `json:"id"`
	UserID              uuid.UUID  `json:"user_id"`
	TotalPrice          float64    `json:"total_price"`
	TotalPriceDiscount  float64    `json:"total_price_discount"`
}