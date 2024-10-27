package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"healthApi/api-gateway/internal/models"
	"healthApi/api-gateway/internal/repositories"
	"healthApi/api-gateway/internal/utils"

	"time"
)

type UserService interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	Authenticate(email, password string) (*models.User, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(user *models.User) error {
	// Check if user already exists
	_, result := s.repo.GetByEmail(user.Email)
	if result.RowsAffected > 0 {
		return utils.NewBadRequestError("Email already registered")
	}
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.NewInternalServerError("Error fetching email")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewInternalServerError("Failed to hash password")
	}
	user.Password = string(hashedPassword)

	// Set default values
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.repo.Create(user)
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}
	return user, nil
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}
	return user, nil
}

func (s *userService) Update(user *models.User) error {
	existing, err := s.GetByID(user.ID)
	if err != nil {
		return err
	}

	// Don't update password through this method
	user.Password = existing.Password
	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}

func (s *userService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *userService) Authenticate(email, password string) (*models.User, error) {
	user, result := s.repo.GetByEmail(email)
	if result.Error != nil {
		return nil, utils.NewUnauthorizedError("Invalid credentials")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, utils.NewUnauthorizedError("Invalid credentials")
	}

	return user, nil
}

func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return utils.NewUnauthorizedError("Invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewInternalServerError("Failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}
