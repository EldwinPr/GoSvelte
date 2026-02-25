package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
)

type BalanceService struct {
	BaseService
	Repo *repositories.CompanyBalanceRepository
}

func (s *BalanceService) GetAllBalance() ([]models.CompanyBalance, error) {
	return s.Repo.FindAll()
}

func (s *BalanceService) GetBalanceByID(id string) (*models.CompanyBalance, error) {
	return s.Repo.FindByID(id)
}
