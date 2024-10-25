package models

import (
	"gorm.io/gorm"
	"time"
)

type Meal struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"userId"`
	Type        string    `gorm:"not null" json:"type"` // breakfast, lunch, dinner, snack
	Description string    `json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
}

type MealRequest struct {
	Type        string    `json:"type" validate:"required,oneof=breakfast lunch dinner snack"`
	Description string    `json:"description" validate:"required,max=500"`
	Date        time.Time `json:"date" validate:"required,ltefield=now"`
}

type MealResponse struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
}
