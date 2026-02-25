package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"time"

	"gorm.io/gorm"
)

type InvoiceService struct {
	BaseService
	InvoiceRepo        *repositories.InvoiceRepository
	DetailRepo         *repositories.InvoiceDetailRepository
	PaymentRepo        *repositories.InvoicePaymentRepository
	TransactionService *TransactionService
}

func (s *InvoiceService) CreateInvoice(invoice *models.Invoice) error {
	return s.InvoiceRepo.Create(invoice)
}

func (s *InvoiceService) EditInvoice(invoice *models.Invoice) error {
	return s.InvoiceRepo.Update(invoice)
}

func (s *InvoiceService) MakeInvoicePayment(payment *models.InvoicePayment) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the payment record
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// 2. Create a Company Transaction for this payment
		transaction := &models.CompanyTransaction{
			Date:          time.Now(),
			Description:   "Payment for Invoice ID: " + payment.InvoiceID,
			Amount:        payment.Amount,
			Type:          "Credit",
			ReferenceID:   payment.ID,
			ReferenceType: "InvoicePayment",
		}
		
		// We'll use the TransactionService logic manually here to keep it in the same DB transaction
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 3. Update Company Balance
		var balance models.CompanyBalance
		if err := tx.First(&balance).Error; err == nil {
			balance.Balance += transaction.Amount
			if err := tx.Save(&balance).Error; err != nil {
				return err
			}
		}

		// 4. Update Invoice Status if fully paid (simple check for prototype)
		var invoice models.Invoice
		if err := tx.Preload("Payments").First(&invoice, "id = ?", payment.InvoiceID).Error; err == nil {
			totalPaid := 0.0
			for _, p := range invoice.Payments {
				totalPaid += p.Amount
			}
			if totalPaid >= invoice.TotalAmount {
				invoice.Status = "Paid"
				if err := tx.Save(&invoice).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *InvoiceService) GetAllInvoices() ([]models.Invoice, error) {
	return s.InvoiceRepo.FindAll()
}

func (s *InvoiceService) GetInvoiceByID(id string) (*models.Invoice, error) {
	return s.InvoiceRepo.FindByID(id)
}
