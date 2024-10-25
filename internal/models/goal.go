package models

import (
	"gorm.io/gorm"
	"time"
)

type Goal struct {
	gorm.Model
	UserID      uint
	Type        string    `gorm:"not null"` // exercise, meal, sleep, hydration
	Target      int       `gorm:"not null"`
	Frequency   string    `gorm:"not null"` // daily, weekly
	StartDate   time.Time `gorm:"not null"`
	EndDate     time.Time `gorm:"not null"`
	Description string
}
