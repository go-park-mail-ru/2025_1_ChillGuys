package dto

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AddressDTO struct {
	ID            uuid.UUID   `json:"id"`
	Label         null.String `json:"label" swaggertype:"primitive,string"`
	Region        null.String `json:"region" swaggertype:"primitive,string"`
	City          null.String `json:"city" swaggertype:"primitive,string"`
	AddressString null.String `json:"addressString" swaggertype:"primitive,string"`
	Coordinate    null.String `json:"coordinate" swaggertype:"primitive,string"`
}

type AddressReqDTO struct {
	Label         null.String `json:"label" swaggertype:"primitive,string"`
	Region        null.String `json:"region" swaggertype:"primitive,string"`
	City          null.String `json:"city" swaggertype:"primitive,string"`
	AddressString null.String `json:"addressString" swaggertype:"primitive,string"`
	Coordinate    null.String `json:"coordinate" swaggertype:"primitive,string"`
}

type GetAddressResDTO struct {
	ID            uuid.UUID   `json:"id" `
	Label         null.String `json:"label" swaggertype:"primitive,string"`
	AddressString null.String `json:"addressString" swaggertype:"primitive,string"`
	Coordinate    null.String `json:"coordinate" swaggertype:"primitive,string"`
}

type GetPointAddressResDTO struct {
	ID            uuid.UUID   `json:"id"`
	AddressString null.String `json:"addressString" swaggertype:"primitive,string"`
	Coordinate    null.String `json:"coordinate" swaggertype:"primitive,string"`
}

type AddressListResponse struct {
	Addresses []GetAddressResDTO `json:"addresses"`
}

type PickupPointListResponse struct {
	PickupPoints []GetPointAddressResDTO `json:"pickupPoints"`
}
