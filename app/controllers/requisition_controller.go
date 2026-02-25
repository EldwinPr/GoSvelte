package controllers

import (
	"encoding/json"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"
	"strconv"

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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	userID := r.URL.Query().Get("user_id")
	pending := r.URL.Query().Get("pending") == "true"

	var result interface{}
	var err error

	if pending {
		result, err = c.Service.GetPendingRequisitions(page, pageSize)
	} else {
		result, err = c.Service.GetPaginatedRequisitions(page, pageSize, userID)
	}

	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, result)
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

	// In a real app, we'd get the current user ID from the context (set by middleware)
	user := r.Context().Value("user").(*models.User)

	if err := c.Service.ApproveRequisition(input.ID, user.ID); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, map[string]string{"message": "Requisition approved"})
}

func (c *RequisitionController) MarkGiven(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// In a real app, we'd get the current user ID from the context (set by middleware)
	user := r.Context().Value("user").(*models.User)

	if err := c.Service.MarkAsGiven(input.ID, user.ID); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, map[string]string{"message": "Requisition marked as given"})
}
