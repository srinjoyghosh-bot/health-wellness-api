package models

import (
	"gorm.io/gorm"
	"time"
)

type Exercise struct {
	gorm.Model
	UserID      uint
	Type        string    `gorm:"not null"`
	Duration    int       `gorm:"not null"` // in minutes
	Intensity   string    `gorm:"not null"` // low, medium, high
	Date        time.Time `gorm:"not null"`
	Description string
}
