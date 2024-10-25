package repositories

import (
	"gorm.io/gorm"
	"healthApi/internal/models"
	"log"
	"time"
)

type ExerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) Create(exercise *models.Exercise) error {
	return r.db.Create(exercise).Error
}

func (r *ExerciseRepository) GetByUserID(userID uint) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.Where("user_id = ?", userID).Find(&exercises).Error
	log.Print("in exercise repo Get by user id", err.Error())
	return exercises, err
}

func (r *ExerciseRepository) GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&exercises).Error
	return exercises, err
}
