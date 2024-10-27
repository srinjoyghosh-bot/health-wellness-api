package services

import (
	"errors"
	"gorm.io/gorm"
	"healthApi/api-gateway/internal/models"
	"healthApi/api-gateway/internal/repositories"
	"healthApi/api-gateway/internal/utils"
	"time"
)

type ExerciseService interface {
	CreateExercise(exercise *models.Exercise) error
	GetByID(id uint) (*models.Exercise, error)
	GetUserExercises(userID uint) ([]models.Exercise, error)
	GetExercisesByDateRange(userID uint, startDate, endDate time.Time) ([]models.Exercise, error)
	Update(exercise *models.Exercise) error
	Delete(id uint) error
	//GetStats(userID uint, startDate, endDate time.Time) (*models.ExerciseStats, error)
}

type exerciseService struct {
	repo repositories.ExerciseRepository
}

func NewExerciseService(repo repositories.ExerciseRepository) ExerciseService {
	return &exerciseService{repo: repo}
}

func (s *exerciseService) CreateExercise(exercise *models.Exercise) error {
	if err := validateExercise(exercise); err != nil {
		return utils.NewBadRequestError(err.Error())
	}

	exercise.CreatedAt = time.Now()
	exercise.UpdatedAt = time.Now()

	return s.repo.Create(exercise)
}

func (s *exerciseService) GetByID(id uint) (*models.Exercise, error) {
	exercise, err := s.repo.GetByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError("Exercise not found")
	}
	return exercise, nil
}

func (s *exerciseService) GetUserExercises(userID uint) ([]models.Exercise, error) {
	exercises, result := s.repo.GetByUserID(userID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return exercises, nil
}

func (s *exerciseService) GetExercisesByDateRange(userID uint, startDate, endDate time.Time) ([]models.Exercise, error) {
	if startDate.After(endDate) {
		return nil, utils.NewBadRequestError("Start date must be before end date")
	}
	return s.repo.GetByDateRange(userID, startDate, endDate)
}

func (s *exerciseService) Update(exercise *models.Exercise) error {
	if err := validateExercise(exercise); err != nil {
		return utils.NewBadRequestError(err.Error())
	}

	existing, err := s.GetByID(exercise.ID)
	if err != nil {
		return err
	}

	// Ensure user owns the exercise
	if existing.UserID != exercise.UserID {
		return utils.NewUnauthorizedError("Not authorized to update this exercise")
	}

	exercise.UpdatedAt = time.Now()
	return s.repo.Update(exercise)
}

func (s *exerciseService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func validateExercise(exercise *models.Exercise) error {
	if exercise.Type == "" {
		return utils.NewBadRequestError("Exercise type is required")
	}

	if exercise.Duration <= 0 {
		return utils.NewBadRequestError("Duration must be greater than 0")
	}

	if !isValidIntensity(exercise.Intensity) {
		return utils.NewBadRequestError("Invalid intensity level")
	}

	if exercise.Date.IsZero() {
		return utils.NewBadRequestError("Date is required")
	}

	if exercise.Date.After(time.Now()) {
		return utils.NewBadRequestError("Exercise date cannot be in the future")
	}

	return nil
}

func isValidIntensity(intensity string) bool {
	validIntensities := map[string]bool{
		"low":      true,
		"moderate": true,
		"high":     true,
	}
	return validIntensities[intensity]
}
