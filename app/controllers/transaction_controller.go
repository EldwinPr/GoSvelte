package controllers

import (
	"encoding/json"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"

	"gorm.io/gorm"
)

type TransactionController struct {
	BaseController
	Service *services.TransactionService
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	transactionRepo := &repositories.CompanyTransactionRepository{BaseRepository: repositories.BaseRepository[models.CompanyTransaction]{DB: db}}
	balanceRepo := &repositories.CompanyBalanceRepository{BaseRepository: repositories.BaseRepository[models.CompanyBalance]{DB: db}}
	
	balanceService := &services.BalanceService{BaseService: services.BaseService{DB: db}, Repo: balanceRepo}
	service := &services.TransactionService{BaseService: services.BaseService{DB: db}, Repo: transactionRepo, BalanceService: balanceService}

	return &TransactionController{Service: service}
}

func (c *TransactionController) Index(w http.ResponseWriter, r *http.Request) {
	transactions, err := c.Service.GetAllTransactions()
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, transactions)
}

func (c *TransactionController) Create(w http.ResponseWriter, r *http.Request) {
	var transaction models.CompanyTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Service.AddNewTransaction(&transaction); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusCreated, transaction)
}
