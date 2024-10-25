package models

import (
	"gorm.io/gorm"
	"time"
)

type Exercise struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"userId" json:"userID"`
	Type        string    `gorm:"not null" json:"type" json:"type"`
	Duration    int       `gorm:"not null" json:"duration"`  // in minutes
	Intensity   string    `gorm:"not null" json:"intensity"` // low, medium, high
	Date        time.Time `gorm:"not null" json:"date"`
	Description string    `json:"description"`
}

type ExerciseRequest struct {
	Type        string    `json:"type" validate:"required,min=2,max=100"`
	Duration    int       `json:"duration" validate:"required,min=1,max=1440"` // max 24 hours
	Intensity   string    `json:"intensity" validate:"required,oneof=low medium high"`
	Date        time.Time `json:"date" validate:"required,ltefield=now"`
	Description string    `json:"description" validate:"omitempty,max=500"`
}

type ExerciseResponse struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Duration    int       `json:"duration"`
	Intensity   string    `json:"intensity"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
