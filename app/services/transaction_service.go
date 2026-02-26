package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TransactionService struct {
	BaseService
	Repo           *repositories.CompanyTransactionRepository
	BalanceService *BalanceService
}

func (s *TransactionService) AddNewTransaction(transaction *models.CompanyTransaction) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Create the transaction record
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// Update Company Balance based on transaction type
		var balance models.CompanyBalance
		if err := tx.First(&balance, "id = ?", transaction.BalanceID).Error; err == nil {
			if transaction.Type == "Credit" {
				balance.Balance += transaction.Amount
			} else {
				balance.Balance -= transaction.Amount
			}
			if err := tx.Save(&balance).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *TransactionService) EditTransaction(transaction *models.CompanyTransaction) error {
	return s.Repo.Update(transaction)
}

func (s *TransactionService) GetMonthlyStats() (income, expense float64, err error) {
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var results []struct {
		Category string
		Total    float64
	}

	err = s.DB.Model(&models.CompanyTransaction{}).
		Where("date >= ?", firstOfMonth).
		Select("category, SUM(amount) as total").
		Group("category").
		Scan(&results).Error

	if err != nil {
		return 0, 0, err
	}

	for _, res := range results {
		if res.Category == "Income" || res.Category == "Sales" {
			income += res.Total
		} else if res.Category == "Expense" || res.Category == "Salary" || res.Category == "Purchase" {
			expense += res.Total
		}
	}

	return income, expense, nil
}

func (s *TransactionService) GetMonthlyTrends(months int) ([]map[string]interface{}, error) {
	now := time.Now()
	// Start from 'months' ago, at the first of the month
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -months+1, 0)

	var results []struct {
		Month  string
		Type   string
		Amount float64
	}

	err := s.DB.Model(&models.CompanyTransaction{}).
		Where("date >= ?", startDate).
		Select("strftime('%Y-%m', date) as month, type, SUM(amount) as amount").
		Group("month, type").
		Order("month ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Prepare map for all months in range
	trendMap := make(map[string]map[string]float64)
	for i := 0; i < months; i++ {
		mStr := startDate.AddDate(0, i, 0).Format("2006-01")
		trendMap[mStr] = map[string]float64{"Credit": 0, "Debit": 0}
	}

	for _, res := range results {
		if _, ok := trendMap[res.Month]; ok {
			trendMap[res.Month][res.Type] = res.Amount
		}
	}

	var finalTrends []map[string]interface{}
	for i := 0; i < months; i++ {
		mStr := startDate.AddDate(0, i, 0).Format("2006-01")
		finalTrends = append(finalTrends, map[string]interface{}{
			"period": mStr,
			"income": trendMap[mStr]["Credit"],
			"out":    trendMap[mStr]["Debit"],
		})
	}

	return finalTrends, nil
}

func (s *TransactionService) GetPaginatedTransactions(page, pageSize int, order string, search string) (*repositories.PaginationResult[models.CompanyTransaction], error) {
	if order == "" {
		order = "company_transactions.date DESC"
	} else if !strings.Contains(order, ".") {
		order = "company_transactions." + order
	}
	
	query := s.DB
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("description LIKE ? OR category LIKE ? OR reference_type LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	return s.Repo.Paginate(page, pageSize, query.Order(order))
}
