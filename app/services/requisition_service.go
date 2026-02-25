package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"strings"
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

func (s *RequisitionService) RejectRequisition(id string, rejectedByID string) error {
	var req models.Requisition
	if err := s.DB.First(&req, "id = ?", id).Error; err != nil {
		return err
	}

	if req.Status != "pending" {
		return errors.New("hanya pengajuan dengan status 'pending' yang bisa ditolak")
	}

	req.Status = "rejected"
	req.ApprovedByID = &rejectedByID // Re-using approved_by field for rejection
	return s.Repo.Update(&req)
}

func (s *RequisitionService) MarkAsGiven(id string, balanceID string, processedByID string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var req models.Requisition
		if err := tx.First(&req, "id = ?", id).Error; err != nil {
			return err
		}

		if req.Status != "approved" {
			return errors.New("hanya pengajuan yang sudah 'approved' yang bisa ditandai 'given'")
		}

		// 1. Check if balance exists
		var balance models.CompanyBalance
		if err := tx.First(&balance, "id = ?", balanceID).Error; err != nil {
			return errors.New("akun keuangan tidak ditemukan")
		}

		if balance.Balance < req.Amount {
			return errors.New("saldo tidak mencukupi di akun terpilih")
		}

		// 2. Update Requisition Status
		req.Status = "given"
		req.ProcessedByID = &processedByID
		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		// 3. Create a Company Transaction (Debit)
		transaction := &models.CompanyTransaction{
			BalanceID:     balanceID,
			Date:          time.Now(),
			Description:   "Pencairan Pengajuan: " + req.Name + " (" + req.Category + ")",
			Category:      "Expense",
			Amount:        req.Amount,
			Type:          "Debit",
			ReferenceID:   req.ID,
			ReferenceType: "Requisition",
		}
		
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 4. Update Company Balance
		balance.Balance -= req.Amount
		if err := tx.Save(&balance).Error; err != nil {
			return err
		}

		// 5. Link transaction to requisition
		req.TransactionID = transaction.ID
		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *RequisitionService) GetPaginatedRequisitions(page, pageSize int, userID string, order string) (*repositories.PaginationResult[models.Requisition], error) {
	if order == "" {
		order = "requisitions.created_at DESC"
	} else if !strings.Contains(order, ".") {
		order = "requisitions." + order
	}
	query := s.DB.Preload("User").Preload("ApprovedBy").Preload("ProcessedBy").Order(order)
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	return s.Repo.Paginate(page, pageSize, query)
}

func (s *RequisitionService) GetPendingRequisitions(page, pageSize int, order string) (*repositories.PaginationResult[models.Requisition], error) {
	if order == "" {
		order = "requisitions.created_at DESC"
	} else if !strings.Contains(order, ".") {
		order = "requisitions." + order
	}
	query := s.DB.Preload("User").Preload("ApprovedBy").Preload("ProcessedBy").Where("status = ?", "pending").Order(order)
	return s.Repo.Paginate(page, pageSize, query)
}

func (s *RequisitionService) GetRequisitionByID(id string) (*models.Requisition, error) {
	var req models.Requisition
	err := s.DB.Preload("User").Preload("ApprovedBy").Preload("ProcessedBy").First(&req, "id = ?", id).Error
	return &req, err
}
