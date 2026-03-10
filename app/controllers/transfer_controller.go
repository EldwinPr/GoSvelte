package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/models"
	"gosvelte/app/services"
	"net/http"
)

type TransferController struct {
	BaseController
	Service *services.TransferService
}

func NewTransferController(service *services.TransferService) *TransferController {
	return &TransferController{
		Service: service,
	}
}

// Create handles the execution of a new transfer.
func (c *TransferController) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var transfer models.Transfer
	if err := c.Bind(r, &transfer); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.ExecuteTransfer(tenantID, &transfer); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, transfer)
}

// List retrieves a paginated list of transfers for the tenant.
func (c *TransferController) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	pagination := c.GetPaginationParams(r)

	result, err := c.Service.GetTransfers(tenantID, pagination)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.PaginatedResponse(w, result)
}

// Delete handles the reversal of a transfer.
func (c *TransferController) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id")

	if err := c.Service.ReverseTransfer(tenantID, id); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusNoContent, nil)
}
