package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AddressDB struct {
	ID            uuid.UUID
	Region        null.String
	City          null.String
	AddressString null.String
	Coordinate    null.String
}

type UserAddress struct {
	ID        uuid.UUID
	Label     null.String
	UserID    uuid.UUID
	AddressID uuid.UUID
}
