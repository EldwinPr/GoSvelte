package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gorm.io/gorm"
)

type FinanceService struct {
	BaseService
	AccountRepo     *repositories.BaseRepository[models.Account]
	TransactionRepo *repositories.BaseRepository[models.Transaction]
}

func NewFinanceService(db *gorm.DB) *FinanceService {
	return &FinanceService{
		BaseService:     BaseService{DB: db},
		AccountRepo:     &repositories.BaseRepository[models.Account]{DB: db},
		TransactionRepo: &repositories.BaseRepository[models.Transaction]{DB: db},
	}
}

// CreateTransaction adds a transaction and updates the account balance atomically.
func (s *FinanceService) CreateTransaction(txData *models.Transaction) error {
	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Get the account (with a lock for safety in high-concurrency)
		var account models.Account
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&account, txData.AccountID).Error; err != nil {
			return err
		}

		// 2. Validate Tenant ownership
		if account.TenantID != txData.TenantID {
			return errors.New("unauthorized account access")
		}

		// 3. Update Balance
		if txData.Type == "Income" {
			account.Balance += txData.Amount
		} else if txData.Type == "Expense" {
			account.Balance -= txData.Amount
		}

		// 4. Save Account
		if err := tx.Save(&account).Error; err != nil {
			return err
		}

		// 5. Create Transaction
		return tx.Create(txData).Error
	})
}

// ReconcileAccount recalculates the balance from the sum of all transactions.
func (s *FinanceService) ReconcileAccount(tenantID, accountID uint) (float64, error) {
	var total float64
	
	// Sum Income
	var income float64
	s.DB.Model(&models.Transaction{}).
		Where("tenant_id = ? AND account_id = ? AND type = ?", tenantID, accountID, "Income").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&income)

	// Sum Expense
	var expense float64
	s.DB.Model(&models.Transaction{}).
		Where("tenant_id = ? AND account_id = ? AND type = ?", tenantID, accountID, "Expense").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&expense)

	total = income - expense

	// Update the account balance to match
	err := s.DB.Model(&models.Account{}).
		Where("id = ? AND tenant_id = ?", accountID, tenantID).
		Update("balance", total).Error

	return total, err
}
