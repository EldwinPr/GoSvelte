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
}

func NewDashboardController(db *gorm.DB) *DashboardController {
	balanceRepo := &repositories.CompanyBalanceRepository{BaseRepository: repositories.BaseRepository[models.CompanyBalance]{DB: db}}
	transactionRepo := &repositories.CompanyTransactionRepository{BaseRepository: repositories.BaseRepository[models.CompanyTransaction]{DB: db}}
	
	balanceService := &services.BalanceService{BaseService: services.BaseService{DB: db}, Repo: balanceRepo}
	transactionService := &services.TransactionService{BaseService: services.BaseService{DB: db}, Repo: transactionRepo, BalanceService: balanceService}

	return &DashboardController{
		BalanceService:     balanceService,
		TransactionService: transactionService,
	}
}

func (c *DashboardController) Stats(w http.ResponseWriter, r *http.Request) {
	balances, _ := c.BalanceService.GetAllBalance()
	// Fetch top 5 recent transactions for stats
	paginatedTransactions, _ := c.TransactionService.GetPaginatedTransactions(1, 5)

	totalBalance := 0.0
	for _, b := range balances {
		totalBalance += b.Balance
	}

	stats := map[string]interface{}{
		"total_balance":       totalBalance,
		"transaction_count":   paginatedTransactions.TotalCount,
		"recent_transactions": paginatedTransactions.Items,
	}

	c.JSON(w, http.StatusOK, stats)
}
