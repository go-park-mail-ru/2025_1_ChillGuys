package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type OrderStatus int

const (
	Placed                    OrderStatus = iota // Оформлен
	InTransit                                    // В пути
)

func (s OrderStatus) String() string {
	return [...]string{
		"placed",
		"in_transit",
	}[s]
}

func ParseOrderStatus(s string) (OrderStatus, error) {
	statuses := [...]string{
		"placed",
		"in_transit",
	}

	for i, val := range statuses {
		if s == val {
			return OrderStatus(i), nil
		}
	}

	return 0, fmt.Errorf("unknown order status: %s", s)
}

func (s OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type OrderPreviewProductDTO struct {
	ProductID       uuid.UUID   `json:"product_id"`
	ProductName     string      `json:"product_name"`
	ProductImageURL null.String `json:"ProductImageURL" swaggertype:"primitive,string"`
	ProductQuantity uint
}
