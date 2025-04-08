package dto

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"time"
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
	ExpectedDeliveryAt *time.Time         `json:"expected_delivery_at"`
	ActualDeliveryAt   *time.Time         `json:"actual_delivery_at"`
	CreatedAt          *time.Time         `json:"created_at,omitempty"`
}

func (orderItem *GetOrderByUserIDResDTO) ConvertToGetOrderByUserIDResDTO(
	address *models.Address,
	products []models.OrderPreviewProduct,
) models.OrderPreview {
	return models.OrderPreview{
		ID:                 orderItem.ID,
		Status:             orderItem.Status,
		TotalPrice:         orderItem.TotalPrice,
		TotalDiscountPrice: orderItem.TotalPriceDiscount,
		Products:           products,
		Address:            *address,
	}
}

type GetOrderProductResDTO struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  uint      `json:"quantity"`
}
