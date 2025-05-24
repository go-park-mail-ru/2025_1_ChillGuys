package models

import (
	"github.com/google/uuid"
	"time"
)

type PromoCode struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Percent    int       `json:"percent"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}