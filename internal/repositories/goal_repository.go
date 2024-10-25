package repositories

import (
	"gorm.io/gorm"
	"healthApi/internal/models"
	"time"
)

type GoalRepository interface {
	Create(goal *models.Goal) error
	GetByID(id uint) (*models.Goal, error)
	GetByUserID(userID uint) ([]models.Goal, error)
	GetActiveGoals(userID uint) ([]models.Goal, error)
	Update(goal *models.Goal) error
	Delete(id uint) error
}

type goalRepository struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) GoalRepository {
	return &goalRepository{db: db}
}

func (r *goalRepository) Create(goal *models.Goal) error {
	return r.db.Create(goal).Error
}

func (r *goalRepository) GetByID(id uint) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.First(&goal, id).Error
	return &goal, err
}

func (r *goalRepository) GetByUserID(userID uint) ([]models.Goal, error) {
	var goals []models.Goal
	err := r.db.Where("user_id = ?", userID).Order("start_date desc").Find(&goals).Error
	return goals, err
}

func (r *goalRepository) GetActiveGoals(userID uint) ([]models.Goal, error) {
	var goals []models.Goal
	now := time.Now()
	err := r.db.Where("user_id = ? AND start_date <= ? AND end_date >= ?", userID, now, now).
		Order("start_date desc").
		Find(&goals).Error
	return goals, err
}

func (r *goalRepository) Update(goal *models.Goal) error {
	return r.db.Save(goal).Error
}

func (r *goalRepository) Delete(id uint) error {
	return r.db.Delete(&models.Goal{}, id).Error
}
