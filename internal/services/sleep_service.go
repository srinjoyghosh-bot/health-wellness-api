package services

import (
	"healthApi/internal/models"
	"healthApi/internal/repositories"
	"healthApi/internal/utils"
	"time"
)

type SleepService interface {
	Create(sleep *models.Sleep) error
	GetByID(id uint) (*models.Sleep, error)
	GetByUserID(userID uint) ([]models.Sleep, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Sleep, error)
	// Update GetWeeklySummary(userID uint) (*models.SleepSummary, error)
	Update(sleep *models.Sleep) error
	Delete(id uint) error
}

type sleepService struct {
	repo repositories.SleepRepository
}

func NewSleepService(repo repositories.SleepRepository) SleepService {
	return &sleepService{repo: repo}
}

func (s *sleepService) Create(sleep *models.Sleep) error {
	if sleep.SleepTime.After(sleep.WakeTime) {
		return utils.NewBadRequestError("Sleep time must be before wake time")
	}

	sleep.Duration = int(sleep.WakeTime.Sub(sleep.SleepTime).Minutes())

	if !isValidSleepQuality(sleep.Quality) {
		return utils.NewBadRequestError("Invalid sleep quality")
	}

	return s.repo.Create(sleep)
}

func (s *sleepService) GetByID(id uint) (*models.Sleep, error) {
	sleep, err := s.repo.GetByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError("Sleep record not found")
	}
	return sleep, nil
}

func (s *sleepService) GetByUserID(userID uint) ([]models.Sleep, error) {
	return s.repo.GetByUserID(userID)
}

func (s *sleepService) GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Sleep, error) {
	if startDate.After(endDate) {
		return nil, utils.NewBadRequestError("Start date must be before end date")
	}
	return s.repo.GetByDateRange(userID, startDate, endDate)
}

func (s *sleepService) Update(sleep *models.Sleep) error {
	if _, err := s.GetByID(sleep.ID); err != nil {
		return err
	}

	sleep.Duration = int(sleep.WakeTime.Sub(sleep.SleepTime).Minutes())
	return s.repo.Update(sleep)
}

func (s *sleepService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func isValidSleepQuality(quality string) bool {
	validQualities := map[string]bool{
		"poor":      true,
		"fair":      true,
		"good":      true,
		"excellent": true,
	}
	return validQualities[quality]
}
