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
	AddressID uuid.UUID            `json:"addressID"`
}

type CreateOrderItemDTO struct {
	ID        uuid.UUID
	ProductID uuid.UUID `json:"productID"`
	Price     float64   `json:"productPrice"`
	Quantity  uint      `json:"quantity"`
}

type CreateOrderRepoReq struct {
	Order             *Order
	UpdatedQuantities map[uuid.UUID]uint
}

type GetOrderByUserIDResDTO struct {
	ID                 uuid.UUID          `json:"id"`
	Status             models.OrderStatus `json:"status"`
	TotalPrice         float64            `json:"totalPrice"`
	TotalPriceDiscount float64            `json:"totalPriceDiscount"`
	AddressID          uuid.UUID          `json:"addressID"`
	ExpectedDeliveryAt *time.Time         `json:"expectedDeliveryAt"`
	ActualDeliveryAt   *time.Time         `json:"actualDeliveryAt"`
	CreatedAt          *time.Time         `json:"createdAt,omitempty"`
}

type OrderPreviewDTO struct {
	ID                 uuid.UUID                       `json:"id"`
	Status             models.OrderStatus              `json:"status"`
	TotalPrice         float64                         `json:"totalPrice"`
	TotalDiscountPrice float64                         `json:"totalDiscountPrice"`
	Products           []models.OrderPreviewProductDTO `json:"products"`
	Address            models.AddressDB                `json:"address"`
	ExpectedDeliveryAt *time.Time                      `json:"expectedDeliveryAt"`
	ActualDeliveryAt   *time.Time                      `json:"actualDeliveryAt"`
	CreatedAt          *time.Time                      `json:"createdAt,omitempty"`
}

func (orderItem *GetOrderByUserIDResDTO) ConvertToGetOrderByUserIDResDTO(
	address *models.AddressDB,
	products []models.OrderPreviewProductDTO,
) OrderPreviewDTO {
	return OrderPreviewDTO{
		ID:                 orderItem.ID,
		Status:             orderItem.Status,
		TotalPrice:         orderItem.TotalPrice,
		TotalDiscountPrice: orderItem.TotalPriceDiscount,
		Products:           products,
		Address:            *address,
	}
}

type GetOrderProductResDTO struct {
	ProductID uuid.UUID `json:"productID"`
	Quantity  uint      `json:"quantity"`
}
