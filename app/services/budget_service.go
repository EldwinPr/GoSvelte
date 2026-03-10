package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"time"

	"gorm.io/gorm"
)

type BudgetService struct {
	BaseService
	Repo *repositories.BaseRepository[models.Budget]
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{
		BaseService: BaseService{DB: db},
		Repo:        &repositories.BaseRepository[models.Budget]{DB: db},
	}
}

// CreateBudget sets a spending limit for a specific Category or a Global limit.
func (s *BudgetService) CreateBudget(tenantID string, budget *models.Budget) error {
	budget.ID = repositories.GenerateULID()
	budget.TenantID = tenantID
	
	if budget.StartDate.IsZero() {
		budget.StartDate = time.Now()
	}
	
	return s.Repo.Create(budget)
}

// GetBudgets retrieves all budgets for a tenant.
func (s *BudgetService) GetBudgets(tenantID string) ([]models.Budget, error) {
	return s.Repo.FindAll(s.Repo.FilterScope("tenant_id", "=", tenantID))
}

// UpdateBudget modifies an existing budget.
func (s *BudgetService) UpdateBudget(tenantID string, budget *models.Budget) error {
	var existing models.Budget
	if err := s.DB.Where("id = ? AND tenant_id = ?", budget.ID, tenantID).First(&existing).Error; err != nil {
		return errors.New("budget not found")
	}

	existing.CategoryID = budget.CategoryID
	existing.Amount = budget.Amount
	existing.Period = budget.Period
	existing.StartDate = budget.StartDate
	existing.EndDate = budget.EndDate
	existing.IsStrict = budget.IsStrict

	return s.Repo.Update(&existing)
}

// DeleteBudget removes a budget.
func (s *BudgetService) DeleteBudget(tenantID string, id string) error {
	var budget models.Budget
	if err := s.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&budget).Error; err != nil {
		return errors.New("budget not found")
	}
	return s.Repo.Delete(&budget)
}

// BudgetPerformance holds the calculated spending vs limit.
type BudgetPerformance struct {
	Budget      models.Budget `json:"budget"`
	ActualSpent float64       `json:"actual_spent"`
	Percentage  float64       `json:"percentage"`
	Remaining   float64       `json:"remaining"`
}

// GetBudgetPerformance returns actual vs. budgeted amounts for the current period.
func (s *BudgetService) GetBudgetPerformance(tenantID string) ([]BudgetPerformance, error) {
	budgets, err := s.GetBudgets(tenantID)
	if err != nil {
		return nil, err
	}

	performance := make([]BudgetPerformance, 0, len(budgets))

	for _, budget := range budgets {
		start, end := s.calculatePeriodBoundaries(budget)
		
		var actual float64
		query := s.DB.Model(&models.Transaction{}).
			Where("tenant_id = ? AND type = ? AND date >= ? AND date <= ?", tenantID, "Expense", start, end)
		
		if budget.CategoryID != nil {
			query = query.Where("category_id = ?", *budget.CategoryID)
		}

		query.Select("COALESCE(SUM(amount), 0)").Scan(&actual)

		percent := 0.0
		if budget.Amount > 0 {
			percent = (actual / budget.Amount) * 100
		}

		performance = append(performance, BudgetPerformance{
			Budget:      budget,
			ActualSpent: actual,
			Percentage:  percent,
			Remaining:   budget.Amount - actual,
		})
	}

	return performance, nil
}

// calculatePeriodBoundaries determines the current start and end dates for a budget period.
func (s *BudgetService) calculatePeriodBoundaries(budget models.Budget) (time.Time, time.Time) {
	now := time.Now()
	var start, end time.Time

	switch budget.Period {
	case "Weekly":
		// Start of current week (Monday)
		weekday := int(now.Weekday())
		if weekday == 0 { weekday = 7 } // Sunday = 7
		start = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	
	case "Monthly":
		// Start of current month
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	
	case "Yearly":
		// Start of current year
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	
	case "One-Time":
		start = budget.StartDate
		if budget.EndDate != nil {
			end = *budget.EndDate
		} else {
			end = time.Date(2099, 12, 31, 23, 59, 59, 0, now.Location())
		}
	
	default: // Default to Monthly
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	}

	return start, end
}
