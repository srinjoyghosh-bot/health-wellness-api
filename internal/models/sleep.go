package models

import (
	"gorm.io/gorm"
	"time"
)

type Sleep struct {
	gorm.Model
	UserID    uint
	SleepTime time.Time `gorm:"not null"`
	WakeTime  time.Time `gorm:"not null"`
	Quality   string    `gorm:"not null"` // poor, fair, good, excellent
	Duration  int       // in minutes
}
