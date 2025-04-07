package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

type User struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	Name        string      `json:"name"`
	Surname     null.String `json:"surname" swaggertype:"primitive,string"`
	ImageURL    null.String `json:"imageURL" swaggertype:"primitive,string"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
}
type UserVersionDB struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Version   int
	UpdatedAt time.Time
}
