package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactionService_RecordTransaction(t *testing.T) {
	tenantID := "tenant_1"

	t.Run("Record Income", func(t *testing.T) {
		db := setupTestDB()
		service := NewTransactionService(db)
		account := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 1000}
		db.Create(&account)

		tx := models.Transaction{
			AccountID: account.ID,
			Amount:    500,
			Type:      "Income",
			Date:      time.Now(),
		}
		err := service.RecordTransaction(tenantID, &tx)
		
		assert.NoError(t, err)
		
		var updatedAccount models.Account
		db.First(&updatedAccount, "id = ?", account.ID)
		assert.Equal(t, 1500.0, updatedAccount.Balance)
	})

	t.Run("Record Expense", func(t *testing.T) {
		db := setupTestDB()
		service := NewTransactionService(db)
		account := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 1000}
		db.Create(&account)

		tx := models.Transaction{
			AccountID: account.ID,
			Amount:    200,
			Type:      "Expense",
			Date:      time.Now(),
		}
		err := service.RecordTransaction(tenantID, &tx)
		
		assert.NoError(t, err)
		
		var updatedAccount models.Account
		db.First(&updatedAccount, "id = ?", account.ID)
		assert.Equal(t, 800.0, updatedAccount.Balance) // 1000 - 200
	})

	t.Run("Unauthorized Account", func(t *testing.T) {
		db := setupTestDB()
		service := NewTransactionService(db)
		account := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 1000}
		db.Create(&account)

		tx := models.Transaction{
			AccountID: account.ID,
			Amount:    100,
			Type:      "Income",
		}
		err := service.RecordTransaction("other_tenant", &tx)
		assert.Error(t, err)
		assert.Equal(t, "account not found or unauthorized", err.Error())
	})
}

func TestTransactionService_UpdateTransaction(t *testing.T) {
	tenantID := "tenant_1"

	t.Run("Change Amount Same Account", func(t *testing.T) {
		db := setupTestDB()
		service := NewTransactionService(db)
		acc1 := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 1000}
		db.Create(&acc1)

		tx := models.Transaction{
			ID: repositories.GenerateULID(),
			TenantID: tenantID,
			AccountID: acc1.ID,
			Amount: 200,
			Type: "Expense",
		}
		db.Create(&tx)
		db.Model(&models.Account{}).Where("id = ?", acc1.ID).Update("balance", 800)

		// Update to 300
		txUpdate := tx
		txUpdate.Amount = 300
		err := service.UpdateTransaction(tenantID, &txUpdate)
		
		assert.NoError(t, err)
		var updatedAcc1 models.Account
		db.First(&updatedAcc1, "id = ?", acc1.ID)
		assert.Equal(t, 700.0, updatedAcc1.Balance) // 1000 - 300
	})

	t.Run("Change Account", func(t *testing.T) {
		db := setupTestDB()
		service := NewTransactionService(db)
		acc1 := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 1000}
		acc2 := models.Account{ID: "acc_2", TenantID: tenantID, Balance: 500}
		db.Create(&acc1)
		db.Create(&acc2)

		tx := models.Transaction{
			ID: repositories.GenerateULID(),
			TenantID: tenantID,
			AccountID: acc1.ID,
			Amount: 100,
			Type: "Expense",
		}
		db.Create(&tx)
		db.Model(&models.Account{}).Where("id = ?", acc1.ID).Update("balance", 900)

		// Move to acc2
		txUpdate := tx
		txUpdate.AccountID = acc2.ID
		err := service.UpdateTransaction(tenantID, &txUpdate)

		assert.NoError(t, err)
		
		var updatedAcc1, updatedAcc2 models.Account
		db.First(&updatedAcc1, "id = ?", acc1.ID)
		db.First(&updatedAcc2, "id = ?", acc2.ID)
		
		assert.Equal(t, 1000.0, updatedAcc1.Balance) // Reverted
		assert.Equal(t, 400.0, updatedAcc2.Balance)  // Deducted from new acc
	})
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	db := setupTestDB()
	service := NewTransactionService(db)
	
	tenantID := "tenant_1"
	account := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 1000}
	db.Create(&account)

	tx := models.Transaction{
		ID: repositories.GenerateULID(),
		TenantID: tenantID,
		AccountID: account.ID,
		Amount: 300,
		Type: "Expense",
	}
	db.Create(&tx)
	db.Model(&models.Account{}).Where("id = ?", account.ID).Update("balance", 700)

	err := service.DeleteTransaction(tenantID, tx.ID)
	assert.NoError(t, err)

	var updatedAccount models.Account
	db.First(&updatedAccount, "id = ?", account.ID)
	assert.Equal(t, 1000.0, updatedAccount.Balance) // Reverted

	var deletedTx models.Transaction
	err = db.First(&deletedTx, "id = ?", tx.ID).Error
	assert.Error(t, err) // Should be soft-deleted (not found by First)
}
