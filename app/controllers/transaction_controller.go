package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/models"
	"gosvelte/app/services"
	"net/http"
)

type TransactionController struct {
	BaseController
	Service *services.TransactionService
}

func NewTransactionController(service *services.TransactionService) *TransactionController {
	return &TransactionController{
		Service: service,
	}
}

// List retrieves a paginated list of transactions (the ledger).
func (c *TransactionController) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	pagination := c.GetPaginationParams(r)

	// Additional filtering can be added here (e.g., filter by account_id, category_id, or date range)
	// For now, we'll just get the full ledger for the tenant.
	result, err := c.Service.GetLedger(tenantID, pagination)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.PaginatedResponse(w, result)
}

// Create records a new transaction and updates the account balance.
func (c *TransactionController) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var tx models.Transaction
	if err := c.Bind(r, &tx); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.RecordTransaction(tenantID, &tx); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, tx)
}

// Update modifies an existing transaction and adjusts account balance accordingly.
func (c *TransactionController) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var tx models.Transaction
	if err := c.Bind(r, &tx); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if tx.ID == "" {
		c.Error(w, http.StatusBadRequest, "transaction id is required")
		return
	}

	if err := c.Service.UpdateTransaction(tenantID, &tx); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, tx)
}

// Delete removes a transaction and reverts its balance impact.
func (c *TransactionController) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id")

	if err := c.Service.DeleteTransaction(tenantID, id); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusNoContent, nil)
}
