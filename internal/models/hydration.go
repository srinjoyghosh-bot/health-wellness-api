package models

import (
	"gorm.io/gorm"
	"time"
)

type Hydration struct {
	gorm.Model
	UserID uint
	Amount int       `gorm:"not null"` // in ml
	Date   time.Time `gorm:"not null"`
}
