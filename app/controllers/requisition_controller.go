package controllers

import (
	"encoding/json"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"

	"gorm.io/gorm"
)

type RequisitionController struct {
	BaseController
	Service *services.RequisitionService
}

func NewRequisitionController(db *gorm.DB) *RequisitionController {
	requisitionRepo := &repositories.RequisitionRepository{BaseRepository: repositories.BaseRepository[models.Requisition]{DB: db}}
	
	transactionRepo := &repositories.CompanyTransactionRepository{BaseRepository: repositories.BaseRepository[models.CompanyTransaction]{DB: db}}
	balanceRepo := &repositories.CompanyBalanceRepository{BaseRepository: repositories.BaseRepository[models.CompanyBalance]{DB: db}}
	
	balanceService := &services.BalanceService{BaseService: services.BaseService{DB: db}, Repo: balanceRepo}
	transactionService := &services.TransactionService{BaseService: services.BaseService{DB: db}, Repo: transactionRepo, BalanceService: balanceService}

	service := &services.RequisitionService{
		BaseService:        services.BaseService{DB: db},
		Repo:               requisitionRepo,
		TransactionService: transactionService,
	}
	return &RequisitionController{Service: service}
}

func (c *RequisitionController) Index(w http.ResponseWriter, r *http.Request) {
	opts := c.ParseQuery(r)
	orderStr := opts.GetOrder("created_at desc")

	userID := r.URL.Query().Get("user_id")
	pending := r.URL.Query().Get("pending") == "true"

	var result interface{}
	var err error

	if pending {
		result, err = c.Service.GetPendingRequisitions(opts.Page, opts.PageSize, orderStr)
	} else {
		result, err = c.Service.GetPaginatedRequisitions(opts.Page, opts.PageSize, userID, orderStr)
	}

	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, result)
}

func (c *RequisitionController) Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		c.Error(w, http.StatusBadRequest, "Missing ID")
		return
	}

	req, err := c.Service.GetRequisitionByID(id)
	if err != nil {
		c.Error(w, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, req)
}

func (c *RequisitionController) Create(w http.ResponseWriter, r *http.Request) {
	var requisition models.Requisition
	if err := json.NewDecoder(r.Body).Decode(&requisition); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Service.MakeNewRequisition(&requisition); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusCreated, requisition)
}

func (c *RequisitionController) Approve(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user := r.Context().Value("user").(*models.User)

	if err := c.Service.ApproveRequisition(input.ID, user.ID); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, map[string]string{"message": "Requisition approved"})
}

func (c *RequisitionController) Reject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user := r.Context().Value("user").(*models.User)

	if err := c.Service.RejectRequisition(input.ID, user.ID); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, map[string]string{"message": "Requisition rejected"})
}

func (c *RequisitionController) MarkGiven(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID        string `json:"id"`
		BalanceID string `json:"balance_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.BalanceID == "" {
		c.Error(w, http.StatusBadRequest, "Silakan pilih akun keuangan")
		return
	}

	user := r.Context().Value("user").(*models.User)

	if err := c.Service.MarkAsGiven(input.ID, input.BalanceID, user.ID); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, map[string]string{"message": "Requisition marked as given"})
}
