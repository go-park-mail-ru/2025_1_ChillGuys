package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ReviewDB struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment" db:"comment"`
}

type Review struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Surname  null.String    `json:"surname"`
	ImageURL null.String    `json:"image_url"`
	Rating   int       `json:"rating"`
	Comment  string    `json:"comment"`
}