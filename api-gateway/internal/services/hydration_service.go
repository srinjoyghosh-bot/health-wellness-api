package services

import (
	"api-gateway/internal/models"
	"api-gateway/internal/repositories"
	"api-gateway/internal/utils"
	"time"
)

type HydrationService interface {
	Create(hydration *models.Hydration) error
	GetByID(id uint) (*models.Hydration, error)
	GetByUserID(userID uint) ([]models.Hydration, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Hydration, error)
	// Update GetDailySummary(userID uint, date time.Time) (*models.HydrationSummary, error)
	Update(hydration *models.Hydration) error
	Delete(id uint) error
}

type hydrationService struct {
	repo repositories.HydrationRepository
}

func NewHydrationService(repo repositories.HydrationRepository) HydrationService {
	return &hydrationService{repo: repo}
}

func (s *hydrationService) Create(hydration *models.Hydration) error {
	if hydration.Amount <= 0 {
		return utils.NewBadRequestError("Amount must be greater than 0")
	}

	if hydration.Amount > 5000 { // 5 liters max per entry
		return utils.NewBadRequestError("Amount exceeds maximum allowed")
	}

	return s.repo.Create(hydration)
}

func (s *hydrationService) GetByID(id uint) (*models.Hydration, error) {
	hydration, err := s.repo.GetByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError("Hydration record not found")
	}
	return hydration, nil
}

func (s *hydrationService) GetByUserID(userID uint) ([]models.Hydration, error) {
	return s.repo.GetByUserID(userID)
}

func (s *hydrationService) GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Hydration, error) {
	if startDate.After(endDate) {
		return nil, utils.NewBadRequestError("Start date must be before end date")
	}
	return s.repo.GetByDateRange(userID, startDate, endDate)
}

func (s *hydrationService) Update(hydration *models.Hydration) error {
	if _, err := s.GetByID(hydration.ID); err != nil {
		return err
	}
	return s.repo.Update(hydration)
}

func (s *hydrationService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
