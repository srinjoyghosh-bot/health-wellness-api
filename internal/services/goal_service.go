package services

import (
	"healthApi/internal/models"
	"healthApi/internal/repositories"
	"healthApi/internal/utils"
	"time"
)

type GoalService interface {
	Create(goal *models.Goal) error
	GetByID(id uint) (*models.Goal, error)
	GetByUserID(userID uint) ([]models.Goal, error)
	GetActiveGoals(userID uint) ([]models.Goal, error)
	// Update CheckGoalProgress(goalID uint) (*models.GoalProgress, error)
	Update(goal *models.Goal) error
	Delete(id uint) error
}

type goalService struct {
	repo repositories.GoalRepository
}

func NewGoalService(repo repositories.GoalRepository) GoalService {
	return &goalService{repo: repo}
}

func (s *goalService) Create(goal *models.Goal) error {
	if !isValidGoalType(goal.Type) {
		return utils.NewBadRequestError("Invalid goal type")
	}

	if !isValidFrequency(goal.Frequency) {
		return utils.NewBadRequestError("Invalid frequency")
	}

	if goal.StartDate.After(goal.EndDate) {
		return utils.NewBadRequestError("Start date must be before end date")
	}

	if goal.Target <= 0 {
		return utils.NewBadRequestError("Target must be greater than 0")
	}

	return s.repo.Create(goal)
}

func (s *goalService) GetByID(id uint) (*models.Goal, error) {
	goal, err := s.repo.GetByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError("Goal not found")
	}
	return goal, nil
}

func (s *goalService) GetByUserID(userID uint) ([]models.Goal, error) {
	return s.repo.GetByUserID(userID)
}

func (s *goalService) GetActiveGoals(userID uint) ([]models.Goal, error) {
	now := time.Now()
	goals, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var activeGoals []models.Goal
	for _, goal := range goals {
		if !goal.StartDate.After(now) && !goal.EndDate.Before(now) {
			activeGoals = append(activeGoals, goal)
		}
	}

	return activeGoals, nil
}

func (s *goalService) Update(goal *models.Goal) error {
	if _, err := s.GetByID(goal.ID); err != nil {
		return err
	}
	return s.repo.Update(goal)
}

func (s *goalService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func isValidGoalType(goalType string) bool {
	validTypes := map[string]bool{
		"exercise":  true,
		"meal":      true,
		"sleep":     true,
		"hydration": true,
	}
	return validTypes[goalType]
}

func isValidFrequency(frequency string) bool {
	validFrequencies := map[string]bool{
		"daily":  true,
		"weekly": true,
	}
	return validFrequencies[frequency]
}
