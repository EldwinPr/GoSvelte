package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
)

type UserService struct {
	BaseService
	Repo *repositories.UserRepository
}

func (s *UserService) GetPaginatedUsers(page, pageSize int, order string, search string) (*repositories.PaginationResult[models.User], error) {
	if order == "" {
		order = "name ASC"
	}
	query := s.DB
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", searchTerm, searchTerm)
	}
	return s.Repo.Paginate(page, pageSize, query.Order(order))
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.FindAll()
}
