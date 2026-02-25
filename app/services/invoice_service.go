package services

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"strings"
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
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Calculate total amount from details if not set
		var total float64
		for i := range invoice.Details {
			invoice.Details[i].Subtotal = float64(invoice.Details[i].Quantity) * invoice.Details[i].UnitPrice
			total += invoice.Details[i].Subtotal
		}
		invoice.TotalAmount = total
		invoice.Status = "Draft"

		if err := tx.Create(invoice).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *InvoiceService) EditInvoice(invoice *models.Invoice) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Calculate new total
		var total float64
		for i := range invoice.Details {
			invoice.Details[i].InvoiceID = invoice.ID
			invoice.Details[i].Subtotal = float64(invoice.Details[i].Quantity) * invoice.Details[i].UnitPrice
			total += invoice.Details[i].Subtotal
		}
		invoice.TotalAmount = total

		// Delete old details and save new ones
		if err := tx.Where("invoice_id = ?", invoice.ID).Delete(&models.InvoiceDetail{}).Error; err != nil {
			return err
		}

		if err := tx.Save(invoice).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *InvoiceService) MakeInvoicePayment(payment *models.InvoicePayment) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the payment record
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// 2. Create a Company Transaction for this payment
		transaction := &models.CompanyTransaction{
			BalanceID:     payment.BalanceID,
			Date:          time.Now(),
			Description:   "Payment for Invoice: " + payment.InvoiceID,
			Category:      "Income",
			Amount:        payment.Amount,
			Type:          "Credit",
			ReferenceID:   payment.InvoiceID, // Link directly to invoice
			ReferenceType: "InvoicePayment",
		}
		
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 3. Update Company Balance
		var balance models.CompanyBalance
		if err := tx.First(&balance, "id = ?", payment.BalanceID).Error; err != nil {
			return err
		}
		balance.Balance += payment.Amount
		if err := tx.Save(&balance).Error; err != nil {
			return err
		}

		// 4. Update Invoice Status if fully paid
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

func (s *InvoiceService) GetPaginatedInvoices(page, pageSize int, order string, search string) (*repositories.PaginationResult[models.Invoice], error) {
	if order == "" {
		order = "invoices.created_at DESC"
	} else if !strings.Contains(order, ".") {
		// Only qualify if not already qualified
		order = "invoices." + order
	}
	
	query := s.DB.Preload("Customer").Preload("Payments")
	
	if search != "" {
		// Search in Invoice Number or Customer Name (via Join)
		searchTerm := "%" + search + "%"
		query = query.Joins("LEFT JOIN customer_credits ON customer_credits.id = invoices.customer_id").
			Where("invoices.number LIKE ? OR customer_credits.customer_name LIKE ?", searchTerm, searchTerm)
	}

	return s.InvoiceRepo.Paginate(page, pageSize, query.Order(order))
}

func (s *InvoiceService) GetPaymentSummary() (float64, error) {
	var summary struct {
		TotalDue float64
	}
	// Use COALESCE to handle empty tables (NULL result from SUM)
	err := s.DB.Model(&models.Invoice{}).
		Where("status != ?", "Paid").
		Select("COALESCE(SUM(total_amount), 0) as total_due").
		Scan(&summary).Error
	
	if err != nil {
		return 0, err
	}

	var totalPaid float64
	err = s.DB.Model(&models.InvoicePayment{}).
		Joins("JOIN invoices ON invoices.id = invoice_payments.invoice_id").
		Where("invoices.status != ?", "Paid").
		Select("COALESCE(SUM(invoice_payments.amount), 0)").
		Scan(&totalPaid).Error

	return summary.TotalDue - totalPaid, err
}

func (s *InvoiceService) GetInvoiceByID(id string) (*models.Invoice, error) {
	var invoice models.Invoice
	// 1. Try direct invoice lookup
	err := s.DB.Preload("Customer").Preload("Details").Preload("Payments").First(&invoice, "id = ?", id).Error
	
	if err != nil && err == gorm.ErrRecordNotFound {
		// 2. If not found, check if the ID is actually a Payment ID
		var payment models.InvoicePayment
		if pErr := s.DB.Select("invoice_id").First(&payment, "id = ?", id).Error; pErr == nil {
			// Found a payment! Use its parent InvoiceID
			return s.GetInvoiceByID(payment.InvoiceID)
		}
	}
	
	return &invoice, err
}

func (s *InvoiceService) GetPaginatedPayments(page, pageSize int, order string) (*repositories.PaginationResult[models.InvoicePayment], error) {
	if order == "" {
		order = "invoice_payments.payment_date DESC"
	} else if !strings.Contains(order, ".") {
		order = "invoice_payments." + order
	}
	return s.PaymentRepo.Paginate(page, pageSize, s.DB.Preload("CreatedBy").Order(order))
}
