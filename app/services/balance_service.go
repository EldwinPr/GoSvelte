package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
)

type BalanceService struct {
	BaseService
	Repo *repositories.CompanyBalanceRepository
}

func (s *BalanceService) GetPaginatedBalance(page, pageSize int, order string, search string) (*repositories.PaginationResult[models.CompanyBalance], error) {
	if order == "" {
		order = "account_name ASC"
	}
	query := s.DB
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("account_name LIKE ?", searchTerm)
	}
	return s.Repo.Paginate(page, pageSize, query.Order(order))
}

func (s *BalanceService) GetAllBalance() ([]models.CompanyBalance, error) {
	return s.Repo.FindAll()
}

func (s *BalanceService) GetBalanceByID(id string) (*models.CompanyBalance, error) {
	return s.Repo.FindByID(id)
}
