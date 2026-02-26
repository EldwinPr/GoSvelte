package controllers

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"

	"gorm.io/gorm"
)

type DashboardController struct {
	BaseController
	BalanceService     *services.BalanceService
	TransactionService *services.TransactionService
	InvoiceService     *services.InvoiceService
	RequisitionService *services.RequisitionService
}

func NewDashboardController(db *gorm.DB) *DashboardController {
	balanceRepo := &repositories.CompanyBalanceRepository{BaseRepository: repositories.BaseRepository[models.CompanyBalance]{DB: db}}
	transactionRepo := &repositories.CompanyTransactionRepository{BaseRepository: repositories.BaseRepository[models.CompanyTransaction]{DB: db}}
	invoiceRepo := &repositories.InvoiceRepository{BaseRepository: repositories.BaseRepository[models.Invoice]{DB: db}}
	requisitionRepo := &repositories.RequisitionRepository{BaseRepository: repositories.BaseRepository[models.Requisition]{DB: db}}
	detailRepo := &repositories.InvoiceDetailRepository{BaseRepository: repositories.BaseRepository[models.InvoiceDetail]{DB: db}}
	paymentRepo := &repositories.InvoicePaymentRepository{BaseRepository: repositories.BaseRepository[models.InvoicePayment]{DB: db}}
	
	balanceService := &services.BalanceService{BaseService: services.BaseService{DB: db}, Repo: balanceRepo}
	transactionService := &services.TransactionService{BaseService: services.BaseService{DB: db}, Repo: transactionRepo, BalanceService: balanceService}
	invoiceService := &services.InvoiceService{
		BaseService: services.BaseService{DB: db},
		InvoiceRepo: invoiceRepo,
		DetailRepo:  detailRepo,
		PaymentRepo: paymentRepo,
		TransactionService: transactionService,
	}
	requisitionService := &services.RequisitionService{
		BaseService: services.BaseService{DB: db},
		Repo:        requisitionRepo,
		TransactionService: transactionService,
	}

	return &DashboardController{
		BalanceService:     balanceService,
		TransactionService: transactionService,
		InvoiceService:     invoiceService,
		RequisitionService: requisitionService,
	}
}

func (c *DashboardController) Stats(w http.ResponseWriter, r *http.Request) {
	balances, _ := c.BalanceService.GetAllBalance()
	paginatedTransactions, _ := c.TransactionService.GetPaginatedTransactions(1, 5, "date DESC", "")
	
	// Get Pending Requisitions
	pendingReqs, _ := c.RequisitionService.GetPendingRequisitions(1, 5, "created_at DESC")
	
	// Get Recent Invoices
	recentInvoices, _ := c.InvoiceService.GetPaginatedInvoices(1, 5, "created_at DESC", "")

	// Get total unpaid amount
	unpaidAmount, _ := c.InvoiceService.GetPaymentSummary()

	// New Analytics
	mtdIncome, mtdExpense, _ := c.TransactionService.GetMonthlyStats()
	monthlyTrends, _ := c.TransactionService.GetMonthlyTrends(6)

	totalBalance := 0.0
	for _, b := range balances {
		totalBalance += b.Balance
	}

	stats := map[string]interface{}{
		"total_balance":        totalBalance,
		"transaction_count":    paginatedTransactions.TotalCount,
		"recent_transactions":  paginatedTransactions.Items,
		"pending_req_count":    pendingReqs.TotalCount,
		"recent_requisitions":  pendingReqs.Items,
		"recent_invoices":      recentInvoices.Items,
		"unpaid_invoice_total": unpaidAmount,
		"mtd_omset":            mtdIncome,
		"mtd_expenses":         mtdExpense,
		"monthly_trends":       monthlyTrends,
	}

	c.JSON(w, http.StatusOK, stats)
}
