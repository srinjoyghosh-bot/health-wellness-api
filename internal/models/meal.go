package models

import (
	"gorm.io/gorm"
	"time"
)

type Meal struct {
	gorm.Model
	UserID      uint
	Type        string `gorm:"not null"` // breakfast, lunch, dinner, snack
	Description string
	Date        time.Time `gorm:"not null"`
}
