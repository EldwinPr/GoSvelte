package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"

	"gorm.io/gorm"
)

type TransactionService struct {
	BaseService
	Repo        *repositories.TransactionRepository
	AccountRepo *repositories.AccountRepository
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		BaseService: BaseService{DB: db},
		Repo: &repositories.TransactionRepository{
			BaseRepository: repositories.BaseRepository[models.Transaction]{DB: db},
		},
		AccountRepo: &repositories.AccountRepository{
			BaseRepository: repositories.BaseRepository[models.Account]{DB: db},
		},
	}
}

// RecordTransaction creates a new transaction and updates account balance atomically.
func (s *TransactionService) RecordTransaction(tenantID string, txData *models.Transaction) error {
	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch and Lock the Account
		var account models.Account
		query := tx.Where("id = ? AND tenant_id = ?", txData.AccountID, tenantID)
		if tx.Dialector.Name() != "sqlite" {
			query = query.Clauses(gorm.Expr("FOR UPDATE"))
		}
		if err := query.First(&account).Error; err != nil {
			return errors.New("account not found or unauthorized")
		}

		// 2. Set Transaction IDs
		txData.ID = repositories.GenerateULID()
		txData.TenantID = tenantID

		// 3. Update Balance
		if txData.Type == "Income" {
			account.Balance += txData.Amount
		} else if txData.Type == "Expense" {
			account.Balance -= txData.Amount
		}

		// 4. Save Account and Create Transaction
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
		return tx.Create(txData).Error
	})
}

// UpdateTransaction handles adjustment of balance based on the delta change.
func (s *TransactionService) UpdateTransaction(tenantID string, txData *models.Transaction) error {
	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch original transaction
		var original models.Transaction
		if err := tx.First(&original, "id = ? AND tenant_id = ?", txData.ID, tenantID).Error; err != nil {
			return errors.New("transaction not found or unauthorized")
		}

		// 2. Handle Account Changes
		if original.AccountID == txData.AccountID {
			// Same account: simple adjustment
			var account models.Account
			query := tx.Where("id = ? AND tenant_id = ?", original.AccountID, tenantID)
			if tx.Dialector.Name() != "sqlite" {
				query = query.Clauses(gorm.Expr("FOR UPDATE"))
			}
			if err := query.First(&account).Error; err != nil {
				return errors.New("account not found")
			}

			// Revert Original Impact
			if original.Type == "Income" {
				account.Balance -= original.Amount
			} else {
				account.Balance += original.Amount
			}

			// Apply New Impact
			if txData.Type == "Income" {
				account.Balance += txData.Amount
			} else {
				account.Balance -= txData.Amount
			}

			if err := tx.Save(&account).Error; err != nil {
				return err
			}
		} else {
			// Different account: complex adjustment
			var oldAccount, newAccount models.Account
			
			// Lock both accounts
			oldQuery := tx.Where("id = ? AND tenant_id = ?", original.AccountID, tenantID)
			newQuery := tx.Where("id = ? AND tenant_id = ?", txData.AccountID, tenantID)
			if tx.Dialector.Name() != "sqlite" {
				oldQuery = oldQuery.Clauses(gorm.Expr("FOR UPDATE"))
				newQuery = newQuery.Clauses(gorm.Expr("FOR UPDATE"))
			}

			if err := oldQuery.First(&oldAccount).Error; err != nil {
				return errors.New("original account not found")
			}
			if err := newQuery.First(&newAccount).Error; err != nil {
				return errors.New("new account not found or unauthorized")
			}

			// Revert from old account
			if original.Type == "Income" {
				oldAccount.Balance -= original.Amount
			} else {
				oldAccount.Balance += original.Amount
			}

			// Apply to new account
			if txData.Type == "Income" {
				newAccount.Balance += txData.Amount
			} else {
				newAccount.Balance -= txData.Amount
			}

			if err := tx.Save(&oldAccount).Error; err != nil {
				return err
			}
			if err := tx.Save(&newAccount).Error; err != nil {
				return err
			}
		}

		// 3. Update Record
		original.AccountID = txData.AccountID
		original.Amount = txData.Amount
		original.Description = txData.Description
		original.Date = txData.Date
		original.CategoryID = txData.CategoryID
		original.Type = txData.Type

		return tx.Save(&original).Error
	})
}

// DeleteTransaction reverts the balance impact and soft-deletes the record.
func (s *TransactionService) DeleteTransaction(tenantID string, id string) error {
	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch the transaction
		var txRecord models.Transaction
		if err := tx.First(&txRecord, "id = ? AND tenant_id = ?", id, tenantID).Error; err != nil {
			return errors.New("transaction not found or unauthorized")
		}

		// 2. Fetch and Lock the Account
		var account models.Account
		query := tx.Where("id = ? AND tenant_id = ?", txRecord.AccountID, tenantID)
		if tx.Dialector.Name() != "sqlite" {
			query = query.Clauses(gorm.Expr("FOR UPDATE"))
		}
		if err := query.First(&account).Error; err != nil {
			return errors.New("account not found")
		}

		// 3. Revert Balance Impact
		if txRecord.Type == "Income" {
			account.Balance -= txRecord.Amount
		} else {
			account.Balance += txRecord.Amount
		}

		// 4. Save and Delete
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
		return tx.Delete(&txRecord).Error
	})
}

// GetLedger returns a paginated list of transactions for a tenant.
func (s *TransactionService) GetLedger(tenantID string, pagination *repositories.Pagination, filters ...func(*gorm.DB) *gorm.DB) (*repositories.Pagination, error) {
	// Add tenant isolation scope as a mandatory filter
	tenantScope := s.Repo.FilterScope("tenant_id", "=", tenantID)
	allScopes := append([]func(*gorm.DB) *gorm.DB{tenantScope}, filters...)

	return s.Repo.PaginatedFind(pagination, allScopes...)
}
