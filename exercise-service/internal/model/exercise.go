package model

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
