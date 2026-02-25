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
	requisitions, err := c.Service.GetAllRequisitions()
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, requisitions)
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

	if err := c.Service.ApproveRequisition(input.ID); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, map[string]string{"message": "Requisition approved"})
}
