package controllers

import (
	"encoding/json"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"
	"time"

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
	opts := c.ParseQuery(r)
	result, err := c.Service.GetPaginatedInvoices(opts.Page, opts.PageSize, opts.GetOrder("invoices.created_at desc"), opts.Search)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, result)
}

func (c *InvoiceController) Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		c.Error(w, http.StatusBadRequest, "Missing ID")
		return
	}

	invoice, err := c.Service.GetInvoiceByID(id)
	if err != nil {
		c.Error(w, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, invoice)
}

func (c *InvoiceController) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Number      string `json:"number"`
		CustomerID  string `json:"customer_id"`
		InvoiceDate string `json:"invoice_date"`
		DueDate     string `json:"due_date"`
		Details     []struct {
			Description string  `json:"description"`
			Quantity    int     `json:"quantity"`
			UnitPrice   float64 `json:"unit_price"`
		} `json:"details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	invoice := models.Invoice{
		Number:     input.Number,
		CustomerID: input.CustomerID,
	}

	if t, err := time.Parse("2006-01-02", input.InvoiceDate); err == nil {
		invoice.InvoiceDate = t
	}
	if t, err := time.Parse("2006-01-02", input.DueDate); err == nil {
		invoice.DueDate = t
	}

	for _, d := range input.Details {
		invoice.Details = append(invoice.Details, models.InvoiceDetail{
			Description: d.Description,
			Quantity:    d.Quantity,
			UnitPrice:   d.UnitPrice,
		})
	}

	if err := c.Service.CreateInvoice(&invoice); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusCreated, invoice)
}

func (c *InvoiceController) Update(w http.ResponseWriter, r *http.Request) {
	var invoice models.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Service.EditInvoice(&invoice); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, invoice)
}

func (c *InvoiceController) Pay(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InvoiceID   string  `json:"invoice_id"`
		BalanceID   string  `json:"balance_id"`
		Amount      float64 `json:"amount"`
		Method      string  `json:"method"`
		PaymentDate string  `json:"payment_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user := r.Context().Value("user").(*models.User)

	payment := models.InvoicePayment{
		InvoiceID:   input.InvoiceID,
		BalanceID:   input.BalanceID,
		CreatedByID: user.ID,
		Amount:      input.Amount,
		Method:      input.Method,
	}

	if t, err := time.Parse("2006-01-02", input.PaymentDate); err == nil {
		payment.PaymentDate = t
	} else {
		payment.PaymentDate = time.Now()
	}

	if err := c.Service.MakeInvoicePayment(&payment); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, payment)
}

func (c *InvoiceController) Payments(w http.ResponseWriter, r *http.Request) {
	opts := c.ParseQuery(r)
	result, err := c.Service.GetPaginatedPayments(opts.Page, opts.PageSize, opts.GetOrder("payment_date desc"))
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, result)
}

func (c *InvoiceController) PaymentSummary(w http.ResponseWriter, r *http.Request) {
	totalDue, err := c.Service.GetPaymentSummary()
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, map[string]float64{"total_due": totalDue})
}
