package dto

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
)

type Order struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Status             models.OrderStatus
	TotalPrice         float64
	TotalPriceDiscount float64
	AddressID          uuid.UUID
	Items              []CreateOrderItemDTO
}

type CreateOrderDTO struct {
	UserID    uuid.UUID
	Items     []CreateOrderItemDTO `json:"items"`
	AddressID uuid.UUID            `json:"address_id"`
}

type CreateOrderItemDTO struct {
	ID        uuid.UUID
	ProductID uuid.UUID `json:"product_id"`
	Price     float64   `json:"product_price"`
	Quantity  uint      `json:"quantity"`
}

type CreateOrderRepoReq struct {
	Order             *Order
	UpdatedQuantities map[uuid.UUID]uint
}

type GetOrderByUserIDResDTO struct {
	ID                 uuid.UUID          `json:"id"`
	Status             models.OrderStatus `json:"status"`
	TotalPrice         float64            `json:"total_price"`
	TotalPriceDiscount float64            `json:"total_price_discount"`
	AddressID          uuid.UUID          `json:"address_id"`
}

type GetOrderProductResDTO struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  uint      `json:"quantity"`
}
