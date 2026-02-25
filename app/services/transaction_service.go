package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"strings"

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
