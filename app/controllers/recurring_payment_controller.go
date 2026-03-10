package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/models"
	"gosvelte/app/services"
	"net/http"
)

type RecurringPaymentController struct {
	BaseController
	Service *services.RecurringPaymentService
}

func NewRecurringPaymentController(service *services.RecurringPaymentService) *RecurringPaymentController {
	return &RecurringPaymentController{
		Service: service,
	}
}

// List retrieves all recurring payment schedules for the tenant.
func (c *RecurringPaymentController) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	schedules, err := c.Service.GetSchedules(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, schedules)
}

// ListPending retrieves all instances waiting for manual confirmation.
func (c *RecurringPaymentController) ListPending(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	instances, err := c.Service.GetPendingInstances(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, instances)
}

// Create sets up a new recurring payment schedule.
func (c *RecurringPaymentController) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var schedule models.RecurringPayment
	if err := c.Bind(r, &schedule); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.CreateSchedule(tenantID, &schedule); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, schedule)
}

// Toggle activates or deactivates a schedule.
func (c *RecurringPaymentController) Toggle(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id")
	activeStr := r.URL.Query().Get("active")
	active := activeStr == "true"

	if err := c.Service.ToggleSchedule(tenantID, id, active); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]string{"message": "schedule updated"})
}

// ProcessCheck triggers a check for due schedules and creates 'pending' instances.
func (c *RecurringPaymentController) ProcessCheck(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	count, err := c.Service.ProcessPending(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]interface{}{
		"generated_count": count,
		"message":         "check complete, pending instances created",
	})
}

// ConfirmRequest defines the input for confirming a bill.
type ConfirmRequest struct {
	InstanceID string  `json:"instance_id"`
	Amount     float64 `json:"amount"`
}

// Confirm handles the manual confirmation of a pending instance.
func (c *RecurringPaymentController) Confirm(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var req ConfirmRequest
	if err := c.Bind(r, &req); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.ConfirmInstance(tenantID, req.InstanceID, req.Amount); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]string{"message": "bill confirmed and transaction recorded"})
}

// Skip marks a pending instance as skipped.
func (c *RecurringPaymentController) Skip(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id")

	if err := c.Service.SkipInstance(tenantID, id); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]string{"message": "bill skipped"})
}
