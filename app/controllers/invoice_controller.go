package controllers

import (
	"encoding/json"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"

	"gorm.io/gorm"
)

type InvoiceController struct {
	BaseController
	Service *services.InvoiceService
}

func NewInvoiceController(db *gorm.DB) *InvoiceController {
	invoiceRepo := &repositories.InvoiceRepository{BaseRepository: repositories.BaseRepository[models.Invoice]{DB: db}}
	detailRepo := &repositories.InvoiceDetailRepository{BaseRepository: repositories.BaseRepository[models.InvoiceDetail]{DB: db}}
	paymentRepo := &repositories.InvoicePaymentRepository{BaseRepository: repositories.BaseRepository[models.InvoicePayment]{DB: db}}
	
	transactionRepo := &repositories.CompanyTransactionRepository{BaseRepository: repositories.BaseRepository[models.CompanyTransaction]{DB: db}}
	balanceRepo := &repositories.CompanyBalanceRepository{BaseRepository: repositories.BaseRepository[models.CompanyBalance]{DB: db}}
	
	balanceService := &services.BalanceService{BaseService: services.BaseService{DB: db}, Repo: balanceRepo}
	transactionService := &services.TransactionService{BaseService: services.BaseService{DB: db}, Repo: transactionRepo, BalanceService: balanceService}

	service := &services.InvoiceService{
		BaseService:        services.BaseService{DB: db},
		InvoiceRepo:        invoiceRepo,
		DetailRepo:         detailRepo,
		PaymentRepo:        paymentRepo,
		TransactionService: transactionService,
	}
	return &InvoiceController{Service: service}
}

func (c *InvoiceController) Index(w http.ResponseWriter, r *http.Request) {
	invoices, err := c.Service.GetAllInvoices()
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, invoices)
}

func (c *InvoiceController) Create(w http.ResponseWriter, r *http.Request) {
	var invoice models.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Service.CreateInvoice(&invoice); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusCreated, invoice)
}

func (c *InvoiceController) Pay(w http.ResponseWriter, r *http.Request) {
	var payment models.InvoicePayment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Service.MakeInvoicePayment(&payment); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, payment)
}
