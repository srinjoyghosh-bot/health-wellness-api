package models

import (
	"gorm.io/gorm"
	"time"
)

type Goal struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"userId"`
	Type        string    `gorm:"not null" json:"type"` // exercise, meal, sleep, hydration
	Target      int       `gorm:"not null" json:"target"`
	Frequency   string    `gorm:"not null" json:"frequency"` // daily, weekly
	StartDate   time.Time `gorm:"not null" json:"startDate"`
	EndDate     time.Time `gorm:"not null" json:"endDate"`
	Description string    `json:"description"`
}

type GoalRequest struct {
	Type        string    `json:"type" validate:"required,oneof=exercise meal sleep hydration"`
	Target      int       `json:"target" validate:"required,min=1"`
	Frequency   string    `json:"frequency" validate:"required,oneof=daily weekly"`
	StartDate   time.Time `json:"startDate" validate:"required,ltefield=EndDate"`
	EndDate     time.Time `json:"endDate" validate:"required,gtefield=StartDate"`
	Description string    `json:"description" validate:"omitempty,max=500"`
}

type GoalResponse struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Target      int       `json:"target"`
	Frequency   string    `json:"frequency"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
