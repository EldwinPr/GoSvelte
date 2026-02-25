package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"time"

	"gorm.io/gorm"
)

type RequisitionService struct {
	BaseService
	Repo               *repositories.RequisitionRepository
	TransactionService *TransactionService
}

func (s *RequisitionService) MakeNewRequisition(requisition *models.Requisition) error {
	requisition.Status = "Pending"
	return s.Repo.Create(requisition)
}

func (s *RequisitionService) ApproveRequisition(id string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var req models.Requisition
		if err := tx.First(&req, "id = ?", id).Error; err != nil {
			return err
		}

		if req.Status == "Completed" {
			return nil // Already approved/completed
		}

		// 1. Update Requisition Status
		req.Status = "Completed"
		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		// 2. Create a Company Transaction (Debit since it's a buy request)
		transaction := &models.CompanyTransaction{
			Date:          time.Now(),
			Description:   "Requisition Fulfillment: " + req.Description,
			Amount:        req.Amount,
			Type:          "Debit",
			ReferenceID:   req.ID,
			ReferenceType: "Requisition",
		}
		
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 3. Update Company Balance
		var balance models.CompanyBalance
		if err := tx.First(&balance).Error; err == nil {
			balance.Balance -= transaction.Amount
			if err := tx.Save(&balance).Error; err != nil {
				return err
			}
		}

		// 4. Link transaction to requisition
		req.TransactionID = transaction.ID
		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *RequisitionService) GetAllRequisitions() ([]models.Requisition, error) {
	return s.Repo.FindAll()
}
