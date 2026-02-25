package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	BaseService
	Repo *repositories.UserRepository
}

func (s *AuthService) Register(user *models.User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.Repo.Create(user)
}

func (s *AuthService) Login(email string, password string) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Verify hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	return s.Repo.FindByID(id)
}
