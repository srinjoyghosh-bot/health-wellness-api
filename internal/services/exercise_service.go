package services

import (
	"healthApi/internal/models"
	"healthApi/internal/repositories"
	"time"
)

type ExerciseService struct {
	repo *repositories.ExerciseRepository
}

func NewExerciseService(repo *repositories.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

func (s *ExerciseService) CreateExercise(exercise *models.Exercise) error {
	return s.repo.Create(exercise)
}

func (s *ExerciseService) GetUserExercises(userID uint) ([]models.Exercise, error) {
	return s.repo.GetByUserID(userID)
}

func (s *ExerciseService) GetExercisesByDateRange(userID uint, startDate, endDate time.Time) ([]models.Exercise, error) {
	return s.repo.GetByDateRange(userID, startDate, endDate)
}
