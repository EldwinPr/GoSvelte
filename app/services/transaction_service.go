package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"

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
		if err := tx.First(&balance).Error; err == nil {
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

func (s *TransactionService) GetAllTransactions() ([]models.CompanyTransaction, error) {
	return s.Repo.FindAll()
}
