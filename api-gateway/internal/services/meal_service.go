package services

import (
	"healthApi/api-gateway/internal/models"
	"healthApi/api-gateway/internal/repositories"
	"healthApi/api-gateway/internal/utils"
	"time"
)

type MealService interface {
	Create(meal *models.Meal) error
	GetByID(id uint) (*models.Meal, error)
	GetByUserID(userID uint) ([]models.Meal, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Meal, error)
	Update(meal *models.Meal) error
	Delete(id uint) error
}

type mealService struct {
	repo repositories.MealRepository
}

func NewMealService(repo repositories.MealRepository) MealService {
	return &mealService{repo: repo}
}

func (s *mealService) Create(meal *models.Meal) error {
	if meal.Type == "" {
		return utils.NewBadRequestError("Meal type is required")
	}

	if !isValidMealType(meal.Type) {
		return utils.NewBadRequestError("Invalid meal type")
	}

	return s.repo.Create(meal)
}

func (s *mealService) GetByID(id uint) (*models.Meal, error) {
	meal, err := s.repo.GetByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError("Meal not found")
	}
	return meal, nil
}

func (s *mealService) GetByUserID(userID uint) ([]models.Meal, error) {
	return s.repo.GetByUserID(userID)
}

func (s *mealService) GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Meal, error) {
	if startDate.After(endDate) {
		return nil, utils.NewBadRequestError("Start date must be before end date")
	}
	return s.repo.GetByDateRange(userID, startDate, endDate)
}

func (s *mealService) Update(meal *models.Meal) error {
	if _, err := s.GetByID(meal.ID); err != nil {
		return err
	}
	return s.repo.Update(meal)
}

func (s *mealService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func isValidMealType(mealType string) bool {
	validTypes := map[string]bool{
		"breakfast": true,
		"lunch":     true,
		"dinner":    true,
		"snack":     true,
	}
	return validTypes[mealType]
}
