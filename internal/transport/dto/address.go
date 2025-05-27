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

type AddressReqDTO struct {
	Label         null.String `json:"label" swaggertype:"primitive,string"`
	Region        null.String `json:"region" swaggertype:"primitive,string"`
	City          null.String `json:"city" swaggertype:"primitive,string"`
	AddressString null.String `json:"addressString" swaggertype:"primitive,string"`
	Coordinate    null.String `json:"coordinate" swaggertype:"primitive,string"`
}
type GeoapifyResponse struct {
	Features []GeoapifyFeature `json:"features"`
}

type GeoapifyFeature struct {
	Properties struct {
		ResultType string  `json:"result_type"`
		Lon        float64 `json:"lon"`
		Lat        float64 `json:"lat"`
		Rank       struct {
			Importance float64 `json:"importance"`
		} `json:"rank"`
	} `json:"properties"`
}
