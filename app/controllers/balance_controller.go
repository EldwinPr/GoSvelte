package controllers

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"

	"gorm.io/gorm"
)

type BalanceController struct {
	BaseController
	Service *services.BalanceService
}

func NewBalanceController(db *gorm.DB) *BalanceController {
	repo := &repositories.CompanyBalanceRepository{BaseRepository: repositories.BaseRepository[models.CompanyBalance]{DB: db}}
	service := &services.BalanceService{BaseService: services.BaseService{DB: db}, Repo: repo}
	return &BalanceController{Service: service}
}

func (c *BalanceController) Index(w http.ResponseWriter, r *http.Request) {
	opts := c.ParseQuery(r)
	result, err := c.Service.GetPaginatedBalance(opts.Page, opts.PageSize, opts.GetOrder("account_name asc"), opts.Search)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, result)
}

func (c *BalanceController) Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id") // Simplified for prototype
	balance, err := c.Service.GetBalanceByID(id)
	if err != nil {
		c.Error(w, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, balance)
}
