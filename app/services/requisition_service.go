package services

import (
	"errors"
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
	requisition.Status = "pending"
	return s.Repo.Create(requisition)
}

func (s *RequisitionService) ApproveRequisition(id string, approvedByID string) error {
	var req models.Requisition
	if err := s.DB.First(&req, "id = ?", id).Error; err != nil {
		return err
	}

	if req.Status != "pending" {
		return errors.New("hanya pengajuan dengan status 'pending' yang bisa disetujui")
	}

	req.Status = "approved"
	req.ApprovedByID = &approvedByID
	return s.Repo.Update(&req)
}

func (s *RequisitionService) MarkAsGiven(id string, processedByID string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var req models.Requisition
		if err := tx.First(&req, "id = ?", id).Error; err != nil {
			return err
		}

		if req.Status != "approved" {
			return errors.New("hanya pengajuan yang sudah 'approved' yang bisa ditandai 'given'")
		}

		// 1. Update Requisition Status
		req.Status = "given"
		req.ProcessedByID = &processedByID
		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		// 2. Create a Company Transaction (Debit)
		transaction := &models.CompanyTransaction{
			Date:          time.Now(),
			Description:   "Pencairan Pengajuan: " + req.Name + " (" + req.Category + ")",
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

func (s *RequisitionService) GetPaginatedRequisitions(page, pageSize int, userID string) (*repositories.PaginationResult[models.Requisition], error) {
	query := s.DB.Order("created_at DESC")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	return s.Repo.Paginate(page, pageSize, query)
}

func (s *RequisitionService) GetPendingRequisitions(page, pageSize int) (*repositories.PaginationResult[models.Requisition], error) {
	query := s.DB.Where("status = ?", "pending").Order("created_at DESC")
	return s.Repo.Paginate(page, pageSize, query)
}
