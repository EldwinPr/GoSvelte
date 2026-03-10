package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"

	"gorm.io/gorm"
)

type TransferService struct {
	BaseService
	Repo *repositories.TransferRepository
}

func NewTransferService(db *gorm.DB) *TransferService {
	return &TransferService{
		BaseService: BaseService{DB: db},
		Repo: &repositories.TransferRepository{
			BaseRepository: repositories.BaseRepository[models.Transfer]{DB: db},
		},
	}
}

// ExecuteTransfer performs an atomic transfer of funds between two accounts.
func (s *TransferService) ExecuteTransfer(tenantID string, transfer *models.Transfer) error {
	if transfer.FromAccountID == transfer.ToAccountID {
		return errors.New("cannot transfer to the same account")
	}

	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch and Lock both accounts
		var fromAccount, toAccount models.Account
		
		fromQuery := tx.Where("id = ? AND tenant_id = ?", transfer.FromAccountID, tenantID)
		toQuery := tx.Where("id = ? AND tenant_id = ?", transfer.ToAccountID, tenantID)
		
		if tx.Dialector.Name() != "sqlite" {
			fromQuery = fromQuery.Clauses(gorm.Expr("FOR UPDATE"))
			toQuery = toQuery.Clauses(gorm.Expr("FOR UPDATE"))
		}

		// Lock source account
		if err := fromQuery.First(&fromAccount).Error; err != nil {
			return errors.New("source account not found or unauthorized")
		}

		// Lock destination account
		if err := toQuery.First(&toAccount).Error; err != nil {
			return errors.New("destination account not found or unauthorized")
		}

		// 2. Set Transfer ID
		transfer.ID = repositories.GenerateULID()
		transfer.TenantID = tenantID

		// 3. Update account balances
		fromAccount.Balance -= transfer.Amount
		toAccount.Balance += transfer.Amount

		// 4. Create Withdrawal Transaction
		withdrawal := &models.Transaction{
			ID:          repositories.GenerateULID(),
			TenantID:    tenantID,
			AccountID:   transfer.FromAccountID,
			TransferID:  &transfer.ID,
			Amount:      transfer.Amount,
			Description: transfer.Description,
			Date:        transfer.Date,
			Type:        "Expense", // Money leaving source account
		}

		// 5. Create Deposit Transaction
		deposit := &models.Transaction{
			ID:          repositories.GenerateULID(),
			TenantID:    tenantID,
			AccountID:   transfer.ToAccountID,
			TransferID:  &transfer.ID,
			Amount:      transfer.Amount,
			Description: transfer.Description,
			Date:        transfer.Date,
			Type:        "Income", // Money entering destination account
		}

		// 6. Save everything
		if err := tx.Save(&fromAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&toAccount).Error; err != nil {
			return err
		}
		if err := tx.Create(withdrawal).Error; err != nil {
			return err
		}
		if err := tx.Create(deposit).Error; err != nil {
			return err
		}
		return tx.Create(transfer).Error
	})
}

// GetTransfers retrieves all transfers for a tenant.
func (s *TransferService) GetTransfers(tenantID string, pagination *repositories.Pagination) (*repositories.Pagination, error) {
	tenantScope := s.Repo.FilterScope("tenant_id", "=", tenantID)
	return s.Repo.PaginatedFind(pagination, tenantScope)
}

// ReverseTransfer handles the rollback of a transfer.
func (s *TransferService) ReverseTransfer(tenantID string, transferID string) error {
	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch the transfer
		var transfer models.Transfer
		if err := tx.First(&transfer, "id = ? AND tenant_id = ?", transferID, tenantID).Error; err != nil {
			return errors.New("transfer not found or unauthorized")
		}

		// 2. Fetch and Lock accounts
		var fromAccount, toAccount models.Account
		fromQuery := tx.Where("id = ? AND tenant_id = ?", transfer.FromAccountID, tenantID)
		toQuery := tx.Where("id = ? AND tenant_id = ?", transfer.ToAccountID, tenantID)
		
		if tx.Dialector.Name() != "sqlite" {
			fromQuery = fromQuery.Clauses(gorm.Expr("FOR UPDATE"))
			toQuery = toQuery.Clauses(gorm.Expr("FOR UPDATE"))
		}

		if err := fromQuery.First(&fromAccount).Error; err != nil {
			return errors.New("source account not found")
		}
		if err := toQuery.First(&toAccount).Error; err != nil {
			return errors.New("destination account not found")
		}

		// 3. Revert balances
		fromAccount.Balance += transfer.Amount
		toAccount.Balance -= transfer.Amount

		// 4. Soft-delete associated transactions
		if err := tx.Where("transfer_id = ? AND tenant_id = ?", transferID, tenantID).Delete(&models.Transaction{}).Error; err != nil {
			return err
		}

		// 5. Save accounts and Soft-delete transfer
		if err := tx.Save(&fromAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&toAccount).Error; err != nil {
			return err
		}
		return tx.Delete(&transfer).Error
	})
}
