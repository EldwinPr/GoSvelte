package services

import (
	"gosvelte/app/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransferService_ExecuteTransfer(t *testing.T) {
	tenantID := "tenant_1"

	t.Run("Successful Transfer", func(t *testing.T) {
		db := setupTestDB()
		service := NewTransferService(db)
		acc1 := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 1000}
		acc2 := models.Account{ID: "acc_2", TenantID: tenantID, Balance: 500}
		db.Create(&acc1)
		db.Create(&acc2)

		transfer := models.Transfer{
			FromAccountID: acc1.ID,
			ToAccountID:   acc2.ID,
			Amount:        200,
			Description:   "Test Transfer",
			Date:          time.Now(),
		}
		
		err := service.ExecuteTransfer(tenantID, &transfer)
		assert.NoError(t, err)

		var updatedAcc1, updatedAcc2 models.Account
		db.First(&updatedAcc1, "id = ?", acc1.ID)
		db.First(&updatedAcc2, "id = ?", acc2.ID)

		assert.Equal(t, 800.0, updatedAcc1.Balance)
		assert.Equal(t, 700.0, updatedAcc2.Balance)

		// Verify 2 transactions were created
		var txCount int64
		db.Model(&models.Transaction{}).Where("transfer_id = ?", transfer.ID).Count(&txCount)
		assert.Equal(t, int64(2), txCount)
	})

	t.Run("Transfer to same account fails", func(t *testing.T) {
		db := setupTestDB()
		service := NewTransferService(db)
		acc1 := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 1000}
		db.Create(&acc1)

		transfer := models.Transfer{
			FromAccountID: acc1.ID,
			ToAccountID:   acc1.ID,
			Amount:        100,
		}
		err := service.ExecuteTransfer(tenantID, &transfer)
		assert.Error(t, err)
		assert.Equal(t, "cannot transfer to the same account", err.Error())
	})
}

func TestTransferService_ReverseTransfer(t *testing.T) {
	db := setupTestDB()
	service := NewTransferService(db)
	
	tenantID := "tenant_1"
	acc1 := models.Account{ID: "acc_1", TenantID: tenantID, Balance: 800}
	acc2 := models.Account{ID: "acc_2", TenantID: tenantID, Balance: 700}
	db.Create(&acc1)
	db.Create(&acc2)

	// Create a transfer and its transactions manually for the test
	transfer := models.Transfer{
		ID:            "transfer_1",
		TenantID:      tenantID,
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        200,
	}
	db.Create(&transfer)

	tx1 := models.Transaction{ID: "tx_1", TenantID: tenantID, TransferID: &transfer.ID, AccountID: acc1.ID, Type: "Expense", Amount: 200}
	tx2 := models.Transaction{ID: "tx_2", TenantID: tenantID, TransferID: &transfer.ID, AccountID: acc2.ID, Type: "Income", Amount: 200}
	db.Create(&tx1)
	db.Create(&tx2)

	t.Run("Reverse Transfer", func(t *testing.T) {
		err := service.ReverseTransfer(tenantID, transfer.ID)
		assert.NoError(t, err)

		var updatedAcc1, updatedAcc2 models.Account
		db.First(&updatedAcc1, "id = ?", acc1.ID)
		db.First(&updatedAcc2, "id = ?", acc2.ID)

		assert.Equal(t, 1000.0, updatedAcc1.Balance) // 800 + 200
		assert.Equal(t, 500.0, updatedAcc2.Balance)  // 700 - 200

		// Verify transactions are deleted
		var txCount int64
		db.Model(&models.Transaction{}).Where("transfer_id = ?", transfer.ID).Count(&txCount)
		assert.Equal(t, int64(0), txCount)
	})
}
