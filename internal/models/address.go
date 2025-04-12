package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AddressDB struct {
	ID            uuid.UUID   `json:"id"`
	Region        null.String `json:"region"`
	City          null.String `json:"city"`
	AddressString null.String `json:"address_string"`
	Coordinate    null.String `json:"coordinate"`
}

type UserAddress struct {
	ID        uuid.UUID   `json:"id"`
	Label     null.String `json:"label"`
	UserID    uuid.UUID   `json:"user_id"`
	AddressID uuid.UUID   `json:"address_id"`
}

type Address struct {
	ID            uuid.UUID   `json:"id"`
	Label         null.String `json:"label"`
	Region        null.String `json:"region"`
	City          null.String `json:"city"`
	AddressString null.String `json:"address_string"`
	Coordinate    null.String `json:"coordinate"`
}

type GetAddressRes struct {
	ID            uuid.UUID   `json:"id"`
	Label         null.String `json:"label"`
	AddressString null.String `json:"address_string"`
	Coordinate    null.String `json:"coordinate"`
}

type GetPointAddressRes struct {
	ID            uuid.UUID   `json:"id"`
	AddressString null.String `json:"address_string"`
	Coordinate    null.String `json:"coordinate"`
}
