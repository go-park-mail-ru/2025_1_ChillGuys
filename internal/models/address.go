package models

import "github.com/guregu/null"

type Address struct {
	City      null.String `json:"city"`
	Street    null.String `json:"street"`
	House     null.String `json:"house"`
	Apartment null.String `json:"apartment"`
	ZipCode   null.String `json:"zipCode"`
}
