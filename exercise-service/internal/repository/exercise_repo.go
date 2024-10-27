package repository

import (
	"exercise-service/internal/model"
	"gorm.io/gorm"
	"log"
	"time"
)

type ExerciseRepository interface {
	Create(exercise *model.Exercise) error
	GetByID(id uint) (*model.Exercise, error)
	GetByUserID(userID uint) ([]model.Exercise, *gorm.DB)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]model.Exercise, error)
	Update(exercise *model.Exercise) error
	Delete(id uint) error
}

type exerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) Create(exercise *model.Exercise) error {
	return r.db.Create(exercise).Error
}

func (r *exerciseRepository) GetByID(id uint) (*model.Exercise, error) {
	var exercise model.Exercise
	err := r.db.First(&exercise, id).Error
	return &exercise, err
}

func (r *exerciseRepository) GetByUserID(userID uint) ([]model.Exercise, *gorm.DB) {
	var exercises []model.Exercise
	result := r.db.Where("user_id = ?", userID).Find(&exercises)
	log.Print("in exercise repo Get by user id", result.Error)
	return exercises, result
}

func (r *exerciseRepository) GetByDateRange(userID uint, startDate, endDate time.Time) ([]model.Exercise, error) {
	var exercises []model.Exercise
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&exercises).Error
	return exercises, err
}

func (r *exerciseRepository) Update(exercise *model.Exercise) error {
	return r.db.Save(exercise).Error
}

func (r *exerciseRepository) Delete(id uint) error {
	return r.db.Delete(&model.Exercise{}, id).Error
}
