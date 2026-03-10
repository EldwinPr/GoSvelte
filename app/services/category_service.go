package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"

	"gorm.io/gorm"
)

type CategoryService struct {
	BaseService
	Repo *repositories.CategoryRepository
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		BaseService: BaseService{DB: db},
		Repo: &repositories.CategoryRepository{
			BaseRepository: repositories.BaseRepository[models.Category]{DB: db},
		},
	}
}

// CreateCategory adds a new income/expense category for a tenant.
func (s *CategoryService) CreateCategory(tenantID string, category *models.Category) error {
	category.ID = repositories.GenerateULID()
	category.TenantID = tenantID
	return s.Repo.Create(category)
}

// GetCategories retrieves all categories for a tenant.
func (s *CategoryService) GetCategories(tenantID string) ([]models.Category, error) {
	return s.Repo.FindAll(s.Repo.FilterScope("tenant_id", "=", tenantID))
}

// UpdateCategory modifies an existing category after verifying ownership.
func (s *CategoryService) UpdateCategory(tenantID string, category *models.Category) error {
	var existing models.Category
	if err := s.DB.Where("id = ? AND tenant_id = ?", category.ID, tenantID).First(&existing).Error; err != nil {
		return errors.New("unauthorized or category not found")
	}

	existing.Name = category.Name
	existing.Icon = category.Icon
	existing.Color = category.Color
	existing.Type = category.Type

	return s.Repo.Update(&existing)
}

// DeleteCategory soft-deletes a category after verifying ownership.
func (s *CategoryService) DeleteCategory(tenantID string, id string) error {
	var category models.Category
	if err := s.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category).Error; err != nil {
		return errors.New("unauthorized or category not found")
	}
	return s.Repo.Delete(&category)
}
