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
	CreatedAt   time.Time `json:"createdAt"`
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

type UpdateGoalRequest struct {
	Type        string    `json:"type,omitempty" validate:"omitempty,oneof=exercise meal sleep hydration"`
	Target      int       `json:"target,omitempty" validate:"omitempty,min=1"`
	Frequency   string    `json:"frequency,omitempty" validate:"omitempty,oneof=daily weekly"`
	StartDate   time.Time `json:"start_date,omitempty" validate:"omitempty,ltefield=EndDate"`
	EndDate     time.Time `json:"end_date,omitempty" validate:"omitempty,gtefield=StartDate"`
	Description string    `json:"description,omitempty" validate:"omitempty,max=500"`
}

func (g *Goal) UpdateFromRequest(req UpdateGoalRequest) {
	// Only update fields that are provided in the request (non-zero values)
	if req.Type != "" {
		g.Type = req.Type
	}

	if req.Target != 0 {
		g.Target = req.Target
	}

	if req.Frequency != "" {
		g.Frequency = req.Frequency
	}

	// For time.Time, we need to check if it's not the zero value
	if !req.StartDate.IsZero() {
		g.StartDate = req.StartDate
	}

	if !req.EndDate.IsZero() {
		g.EndDate = req.EndDate
	}

	// For string fields that allow empty values, you might want to add a specific
	// logic to determine if the field should be updated
	if req.Description != "" {
		g.Description = req.Description
	}
}

func (g *Goal) ToResponse() GoalResponse {
	return GoalResponse{
		ID:          g.ID,
		Type:        g.Type,
		Target:      g.Target,
		Frequency:   g.Frequency,
		StartDate:   g.StartDate,
		EndDate:     g.EndDate,
		Description: g.Description,
		CreatedAt:   g.CreatedAt,
	}
}
