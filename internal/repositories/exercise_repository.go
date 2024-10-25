package repositories

import (
	"gorm.io/gorm"
	"healthApi/internal/models"
	"log"
	"time"
)

type ExerciseRepository interface {
	Create(exercise *models.Exercise) error
	GetByID(id uint) (*models.Exercise, error)
	GetByUserID(userID uint) ([]models.Exercise, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Exercise, error)
	Update(exercise *models.Exercise) error
	Delete(id uint) error
}

type exerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) Create(exercise *models.Exercise) error {
	return r.db.Create(exercise).Error
}

func (r *exerciseRepository) GetByID(id uint) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.db.First(&exercise, id).Error
	return &exercise, err
}

func (r *exerciseRepository) GetByUserID(userID uint) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.Where("user_id = ?", userID).Find(&exercises).Error
	log.Print("in exercise repo Get by user id", err.Error())
	return exercises, err
}

func (r *exerciseRepository) GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&exercises).Error
	return exercises, err
}

func (r *exerciseRepository) Update(exercise *models.Exercise) error {
	return r.db.Save(exercise).Error
}

func (r *exerciseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Exercise{}, id).Error
}
