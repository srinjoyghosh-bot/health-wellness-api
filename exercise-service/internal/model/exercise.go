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

type ExerciseRequest struct {
	Type        string    `json:"type" validate:"required,min=2,max=100"`
	Duration    int       `json:"duration" validate:"required,min=1,max=1440"` // max 24 hours
	Intensity   string    `json:"intensity" validate:"required,oneof=low medium high"`
	Date        time.Time `json:"date" validate:"required"`
	Description string    `json:"description" validate:"omitempty,max=500"`
}

type ExerciseUpdateRequest struct {
	Type        string    `json:"type,omitempty" validate:"omitempty,min=2,max=100"`
	Duration    int       `json:"duration,omitempty" validate:"omitempty,min=1,max=1440"` // max 24 hours
	Intensity   string    `json:"intensity,omitempty" validate:"omitempty,oneof=low medium high"`
	Date        time.Time `json:"date,omitempty" validate:"omitempty"`
	Description string    `json:"description,omitempty" validate:"omitempty,max=500"`
}

type ExerciseResponse struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Duration    int       `json:"duration"`
	Intensity   string    `json:"intensity"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

func (e *Exercise) UpdateFromRequest(req ExerciseUpdateRequest) {
	if req.Type != "" {
		e.Type = req.Type
	}
	if req.Duration != 0 {
		e.Duration = req.Duration
	}
	if req.Intensity != "" {
		e.Intensity = req.Intensity
	}
	if req.Description != "" {
		e.Description = req.Description
	}
	if !req.Date.IsZero() {
		e.Date = req.Date
	}
}

func (e *Exercise) ToResponse() ExerciseResponse {
	return ExerciseResponse{
		ID:          e.ID,
		Type:        e.Type,
		Duration:    e.Duration,
		Intensity:   e.Intensity,
		Date:        e.Date,
		Description: e.Description,
	}
}
