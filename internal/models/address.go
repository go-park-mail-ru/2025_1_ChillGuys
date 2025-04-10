package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AddressDB struct {
	ID        uuid.UUID   `json:"id"`
	City      null.String `json:"city"`
	Street    null.String `json:"street"`
	House     null.String `json:"house"`
	Apartment null.String `json:"apartment"`
	ZipCode   null.String `json:"zip_code"`
	//Coordinate null.String `json:"coordinate"`
}

type UserAddress struct {
	ID        uuid.UUID   `json:"id"`
	Label     null.String `json:"label"`
	UserID    uuid.UUID   `json:"user_id"`
	AddressID uuid.UUID   `json:"address_id"`
}

type Address struct {
	ID        uuid.UUID   `json:"id"`
	Label     null.String `json:"label"`
	City      null.String `json:"city"`
	Street    null.String `json:"street"`
	House     null.String `json:"house"`
	Apartment null.String `json:"apartment"`
	ZipCode   null.String `json:"zip_code"`
}
