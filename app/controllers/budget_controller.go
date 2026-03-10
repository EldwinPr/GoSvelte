package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/models"
	"gosvelte/app/services"
	"net/http"
)

type BudgetController struct {
	BaseController
	Service *services.BudgetService
}

func NewBudgetController(service *services.BudgetService) *BudgetController {
	return &BudgetController{
		Service: service,
	}
}

// List retrieves all budgets for the current tenant.
func (c *BudgetController) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	budgets, err := c.Service.GetBudgets(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, budgets)
}

// Create handles the creation of a new budget.
func (c *BudgetController) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var budget models.Budget
	if err := c.Bind(r, &budget); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.CreateBudget(tenantID, &budget); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, budget)
}

// Update handles budget modifications.
func (c *BudgetController) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var budget models.Budget
	if err := c.Bind(r, &budget); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if budget.ID == "" {
		c.Error(w, http.StatusBadRequest, "budget id is required")
		return
	}

	if err := c.Service.UpdateBudget(tenantID, &budget); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, budget)
}

// Delete handles budget removal.
func (c *BudgetController) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id")

	if err := c.Service.DeleteBudget(tenantID, id); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusNoContent, nil)
}

// Performance returns the spending vs. limit analytics for all budgets.
func (c *BudgetController) Performance(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	results, err := c.Service.GetBudgetPerformance(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, results)
}
