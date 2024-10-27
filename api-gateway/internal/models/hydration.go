package models

import (
	"gorm.io/gorm"
	"time"
)

type Hydration struct {
	gorm.Model
	UserID uint      `gorm:"not null" json:"userId"`
	Amount int       `gorm:"not null" json:"amount"` // in ml
	Date   time.Time `gorm:"not null" json:"date"`
}

type HydrationRequest struct {
	Amount int `json:"amount" validate:"required,min=1,max=5000"` // max 5 liters at once
}

type HydrationResponse struct {
	ID     uint      `json:"id"`
	Amount int       `json:"amount"`
	Date   time.Time `json:"date"`
}

func (h Hydration) ToResponse() HydrationResponse {
	return HydrationResponse{
		ID:     h.ID,
		Amount: h.Amount,
		Date:   h.Date,
	}
}
