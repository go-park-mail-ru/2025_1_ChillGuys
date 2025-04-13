package dto

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AddressDTO struct {
	ID            uuid.UUID   `json:"id"`
	Label         null.String `json:"label"`
	Region        null.String `json:"region"`
	City          null.String `json:"city"`
	AddressString null.String `json:"addressString"`
	Coordinate    null.String `json:"coordinate"`
}

type GetAddressResDTO struct {
	ID            uuid.UUID   `json:"id"`
	Label         null.String `json:"label"`
	AddressString null.String `json:"addressString"`
	Coordinate    null.String `json:"coordinate"`
}

type GetPointAddressResDTO struct {
	ID            uuid.UUID   `json:"id"`
	AddressString null.String `json:"addressString"`
	Coordinate    null.String `json:"coordinate"`
}
