package repositories

import (
	"gorm.io/gorm"
	"healthApi/api-gateway/internal/models"
	"time"
)

type HydrationRepository interface {
	Create(hydration *models.Hydration) error
	GetByID(id uint) (*models.Hydration, error)
	GetByUserID(userID uint) ([]models.Hydration, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Hydration, error)
	Update(hydration *models.Hydration) error
	Delete(id uint) error
}

type hydrationRepository struct {
	db *gorm.DB
}

func NewHydrationRepository(db *gorm.DB) HydrationRepository {
	return &hydrationRepository{db: db}
}

func (r *hydrationRepository) Create(hydration *models.Hydration) error {
	return r.db.Create(hydration).Error
}

func (r *hydrationRepository) GetByID(id uint) (*models.Hydration, error) {
	var hydration models.Hydration
	err := r.db.First(&hydration, id).Error
	return &hydration, err
}

func (r *hydrationRepository) GetByUserID(userID uint) ([]models.Hydration, error) {
	var hydrations []models.Hydration
	err := r.db.Where("user_id = ?", userID).Order("date desc").Find(&hydrations).Error
	return hydrations, err
}

func (r *hydrationRepository) GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Hydration, error) {
	var hydrations []models.Hydration
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date desc").
		Find(&hydrations).Error
	return hydrations, err
}

func (r *hydrationRepository) Update(hydration *models.Hydration) error {
	return r.db.Save(hydration).Error
}

func (r *hydrationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Hydration{}, id).Error
}
