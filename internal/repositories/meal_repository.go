package repositories

import (
	"gorm.io/gorm"
	"healthApi/internal/models"
	"time"
)

type MealRepository interface {
	Create(meal *models.Meal) error
	GetByID(id uint) (*models.Meal, error)
	GetByUserID(userID uint) ([]models.Meal, error)
	GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Meal, error)
	Update(meal *models.Meal) error
	Delete(id uint) error
}

type mealRepository struct {
	db *gorm.DB
}

func NewMealRepository(db *gorm.DB) MealRepository {
	return &mealRepository{db: db}
}

func (r *mealRepository) Create(meal *models.Meal) error {
	return r.db.Create(meal).Error
}

func (r *mealRepository) GetByID(id uint) (*models.Meal, error) {
	var meal models.Meal
	err := r.db.First(&meal, id).Error
	return &meal, err
}

func (r *mealRepository) GetByUserID(userID uint) ([]models.Meal, error) {
	var meals []models.Meal
	err := r.db.Where("user_id = ?", userID).Order("date desc").Find(&meals).Error
	return meals, err
}

func (r *mealRepository) GetByDateRange(userID uint, startDate, endDate time.Time) ([]models.Meal, error) {
	var meals []models.Meal
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date desc").
		Find(&meals).Error
	return meals, err
}

func (r *mealRepository) Update(meal *models.Meal) error {
	return r.db.Save(meal).Error
}

func (r *mealRepository) Delete(id uint) error {
	return r.db.Delete(&models.Meal{}, id).Error
}
