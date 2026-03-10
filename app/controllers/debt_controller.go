package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/models"
	"gosvelte/app/services"
	"net/http"
)

type DebtController struct {
	BaseController
	Service *services.DebtService
}

func NewDebtController(service *services.DebtService) *DebtController {
	return &DebtController{
		Service: service,
	}
}

// List retrieves all debts for the tenant.
func (c *DebtController) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	debts, err := c.Service.GetDebts(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, debts)
}

// Create handles the creation of a new debt (Lent or Borrowed).
func (c *DebtController) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var debt models.Debt
	if err := c.Bind(r, &debt); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.CreateDebt(tenantID, &debt); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, debt)
}

// AddInstallmentRequest defines the input for recording a payment toward a debt.
type AddInstallmentRequest struct {
	Installment models.DebtInstallment `json:"installment"`
	AccountID   string                 `json:"account_id"`
}

// AddInstallment handles the recording of a payment toward a debt.
func (c *DebtController) AddInstallment(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var req AddInstallmentRequest
	if err := c.Bind(r, &req); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.AddInstallment(tenantID, &req.Installment, req.AccountID); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, req.Installment)
}

// Summary retrieves the total lent and borrowed amounts for the tenant.
func (c *DebtController) Summary(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	summary, err := c.Service.GetDebtSummary(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, summary)
}
