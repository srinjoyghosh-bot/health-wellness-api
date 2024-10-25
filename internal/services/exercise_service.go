package services

import (
	"healthApi/internal/models"
	"healthApi/internal/repositories"
	"time"
)

type ExerciseService interface {
	CreateExercise(exercise *models.Exercise) error
	GetByID(id uint) (*models.Exercise, error)
	GetUserExercises(userID uint) ([]models.Exercise, error)
	GetExercisesByDateRange(userID uint, startDate, endDate time.Time) ([]models.Exercise, error)
	Update(exercise *models.Exercise) error
	Delete(id uint) error
}

type exerciseService struct {
	repo repositories.ExerciseRepository
}

func NewExerciseService(repo repositories.ExerciseRepository) ExerciseService {
	return &exerciseService{repo: repo}
}

func (s *exerciseService) CreateExercise(exercise *models.Exercise) error {
	return s.repo.Create(exercise)
}

func (s *exerciseService) GetByID(id uint) (*models.Exercise, error) {
	return s.repo.GetByID(id)
}

func (s *exerciseService) GetUserExercises(userID uint) ([]models.Exercise, error) {
	return s.repo.GetByUserID(userID)
}

func (s *exerciseService) GetExercisesByDateRange(userID uint, startDate, endDate time.Time) ([]models.Exercise, error) {
	return s.repo.GetByDateRange(userID, startDate, endDate)
}

func (s *exerciseService) Update(exercise *models.Exercise) error {
	return s.repo.Update(exercise)
}

func (s *exerciseService) Delete(id uint) error {
	return s.repo.Delete(id)
}
