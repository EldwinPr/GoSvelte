package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/models"
	"gosvelte/app/services"
	"net/http"
)

type AccountController struct {
	BaseController
	Service *services.AccountService
}

func NewAccountController(service *services.AccountService) *AccountController {
	return &AccountController{
		Service: service,
	}
}

// List retrieves all accounts for the current tenant.
func (c *AccountController) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	accounts, err := c.Service.GetAccounts(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, accounts)
}

// Create handles the creation of a new account.
func (c *AccountController) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var account models.Account
	if err := c.Bind(r, &account); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.CreateAccount(tenantID, &account); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, account)
}

// Get retrieves a single account by ID.
func (c *AccountController) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id") // Or use a proper router param if using a more complex router

	account, err := c.Service.GetAccount(tenantID, id)
	if err != nil {
		c.Error(w, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, account)
}

// Update handles account modifications.
func (c *AccountController) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var account models.Account
	if err := c.Bind(r, &account); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// ID must be provided in the body or URL
	if account.ID == "" {
		c.Error(w, http.StatusBadRequest, "account id is required")
		return
	}

	if err := c.Service.UpdateAccount(tenantID, &account); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, account)
}

// Delete handles account removal.
func (c *AccountController) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id")

	if err := c.Service.DeleteAccount(tenantID, id); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusNoContent, nil)
}

// DashboardSummary returns the net worth and breakdown for the tenant.
func (c *AccountController) DashboardSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	summary, err := c.Service.GetDashboardSummary(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, summary)
}
