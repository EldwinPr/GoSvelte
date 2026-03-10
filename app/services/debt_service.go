package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"

	"gorm.io/gorm"
)

type DebtService struct {
	BaseService
	Repo            *repositories.BaseRepository[models.Debt]
	InstallmentRepo *repositories.BaseRepository[models.DebtInstallment]
}

func NewDebtService(db *gorm.DB) *DebtService {
	return &DebtService{
		BaseService: BaseService{DB: db},
		Repo: &repositories.BaseRepository[models.Debt]{DB: db},
		InstallmentRepo: &repositories.BaseRepository[models.DebtInstallment]{DB: db},
	}
}

// CreateDebt records a new Lent or Borrowed fund.
func (s *DebtService) CreateDebt(tenantID string, debt *models.Debt) error {
	debt.ID = repositories.GenerateULID()
	debt.TenantID = tenantID
	debt.Remaining = debt.TotalAmount
	return s.Repo.Create(debt)
}

// GetDebts retrieves all active debts for a tenant.
func (s *DebtService) GetDebts(tenantID string) ([]models.Debt, error) {
	return s.Repo.FindAll(s.Repo.FilterScope("tenant_id", "=", tenantID))
}

// AddInstallment records a payment toward a debt and updates the remaining balance.
func (s *DebtService) AddInstallment(tenantID string, installment *models.DebtInstallment, accountID string) error {
	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch and Lock the Debt
		var debt models.Debt
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&debt, "id = ? AND tenant_id = ?", installment.DebtID, tenantID).Error; err != nil {
			return errors.New("debt not found or unauthorized")
		}

		// 2. Lock and Update the Account
		var account models.Account
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&account, "id = ? AND tenant_id = ?", accountID, tenantID).Error; err != nil {
			return errors.New("account not found")
		}

		// Update balance based on debt type
		// If we are paying a Borrowed debt, it's an Expense
		// If we are receiving payment for a Lent debt, it's an Income
		if debt.Type == "Borrowed" {
			account.Balance -= installment.Amount
		} else {
			account.Balance += installment.Amount
		}
		if err := tx.Save(&account).Error; err != nil {
			return err
		}

		// 3. Create Transaction for the payment
		transaction := &models.Transaction{
			ID:                repositories.GenerateULID(),
			TenantID:          tenantID,
			AccountID:         accountID,
			DebtID:            &debt.ID,
			DebtInstallmentID: &installment.ID,
			Amount:            installment.Amount,
			Description:       "Debt payment: " + debt.Name,
			Date:              installment.PaidDate,
			Type:              "Expense",
		}
		if debt.Type == "Lent" {
			transaction.Type = "Income"
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 4. Update Debt remaining balance
		debt.Remaining -= installment.Amount
		if err := tx.Save(&debt).Error; err != nil {
			return err
		}

		// 5. Create Installment record
		installment.ID = repositories.GenerateULID()
		installment.TenantID = tenantID
		return tx.Create(installment).Error
	})
}

// GetDebtSummary calculates the total lent and borrowed amounts for a tenant.
func (s *DebtService) GetDebtSummary(tenantID string) (map[string]float64, error) {
	var results []struct {
		Type  string
		Total float64
	}
	err := s.DB.Model(&models.Debt{}).
		Select("type, sum(remaining) as total").
		Where("tenant_id = ?", tenantID).
		Group("type").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}

	summary := map[string]float64{
		"Lent":     0,
		"Borrowed": 0,
	}
	for _, res := range results {
		summary[res.Type] = res.Total
	}
	return summary, nil
}
