package models

import (
	"gorm.io/gorm"
	"time"
)

type Sleep struct {
	gorm.Model
	UserID    uint      `gorm:"not null" json:"userId"`
	SleepTime time.Time `gorm:"not null" json:"sleepTime"`
	WakeTime  time.Time `gorm:"not null" json:"wakeTime"`
	Duration  int       `gorm:"not null" json:"duration"`
	Quality   string    `gorm:"not null" json:"quality"` // poor, fair, good, excellent 	// in minutes
	CreatedAt time.Time `sql:"default:CURRENT_TIMESTAMP" json:"createdAt"`
}

type SleepRequest struct {
	SleepTime time.Time `json:"sleepTime" validate:"required,ltefield=WakeTime"`
	WakeTime  time.Time `json:"wakeTime" validate:"required,gtefield=SleepTime"`
	Quality   string    `json:"quality" validate:"required,oneof=poor fair good excellent"`
}

type SleepResponse struct {
	ID        uint      `json:"id"`
	SleepTime time.Time `json:"sleepTime"`
	WakeTime  time.Time `json:"wakeTime"`
	Duration  int       `json:"duration"`
	Quality   string    `json:"quality"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s Sleep) ToResponse() SleepResponse {
	return SleepResponse{
		ID:        s.ID,
		SleepTime: s.SleepTime,
		WakeTime:  s.WakeTime,
		Duration:  s.Duration,
		Quality:   s.Quality,
		CreatedAt: s.CreatedAt,
	}
}
