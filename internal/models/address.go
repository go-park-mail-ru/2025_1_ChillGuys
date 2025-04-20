package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AddressDB struct {
	ID            uuid.UUID   `json:"id"`
	Label         null.String `json:"label" swaggertype:"primitive,string"`
	Region        null.String `json:"region" swaggertype:"primitive,string"`
	City          null.String `json:"city" swaggertype:"primitive,string"`
	AddressString null.String `json:"AddressString" swaggertype:"primitive,string"`
	Coordinate    null.String `json:"coordinate" swaggertype:"primitive,string"`
}

type UserAddress struct {
	ID        uuid.UUID
	Label     null.String `json:"label" swaggertype:"primitive,string"`
	UserID    uuid.UUID
	AddressID uuid.UUID
}
