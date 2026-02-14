package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
)

type UserService struct {
	BaseService
	Repo *repositories.UserRepository
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.FindAll()
}
