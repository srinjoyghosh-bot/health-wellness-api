package repositories

import (
	"gorm.io/gorm"
	"healthApi/internal/models"
	"time"
)

type SleepRepository interface {
	Create(sleep *models.Sleep) error
	GetByID(id uint) (*models.Sleep, error)
	GetByUserID(userID uint) ([]models.Sleep, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Sleep, error)
	Update(sleep *models.Sleep) error
	Delete(id uint) error
}

type sleepRepository struct {
	db *gorm.DB
}

func NewSleepRepository(db *gorm.DB) SleepRepository {
	return &sleepRepository{db: db}
}

func (r *sleepRepository) Create(sleep *models.Sleep) error {
	return r.db.Create(sleep).Error
}

func (r *sleepRepository) GetByID(id uint) (*models.Sleep, error) {
	var sleep models.Sleep
	err := r.db.First(&sleep, id).Error
	return &sleep, err
}

func (r *sleepRepository) GetByUserID(userID uint) ([]models.Sleep, error) {
	var sleeps []models.Sleep
	err := r.db.Where("user_id = ?", userID).Order("sleep_time desc").Find(&sleeps).Error
	return sleeps, err
}

func (r *sleepRepository) GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Sleep, error) {
	var sleeps []models.Sleep
	err := r.db.Where("user_id = ? AND sleep_time BETWEEN ? AND ?", userID, startDate, endDate).
		Order("sleep_time desc").
		Find(&sleeps).Error
	return sleeps, err
}

func (r *sleepRepository) Update(sleep *models.Sleep) error {
	return r.db.Save(sleep).Error
}

func (r *sleepRepository) Delete(id uint) error {
	return r.db.Delete(&models.Sleep{}, id).Error
}
