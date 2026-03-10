package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/services"
	"net/http"
	"strconv"
	"time"
)

type ReportingController struct {
	BaseController
	Service *services.ReportingService
}

func NewReportingController(service *services.ReportingService) *ReportingController {
	return &ReportingController{
		Service: service,
	}
}

// GetDashboard returns a complete snapshot of all dashboard-related analytics.
func (c *ReportingController) GetDashboard(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	snapshot, err := c.Service.GetDashboardSnapshot(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, snapshot)
}

// GetSpendingBreakdown returns expenses grouped by category for a specific date range.
func (c *ReportingController) GetSpendingBreakdown(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	query := r.URL.Query()
	
	// Default to current month if no dates provided
	start, _ := time.Parse("2006-01-02", query.Get("start"))
	end, _ := time.Parse("2006-01-02", query.Get("end"))
	
	if start.IsZero() {
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = now
	}

	results, err := c.Service.GetSpendingBreakdown(tenantID, start, end)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, results)
}

// GetCashFlowHistory returns monthly income vs expense history.
func (c *ReportingController) GetCashFlowHistory(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	monthsStr := r.URL.Query().Get("months")
	months, _ := strconv.Atoi(monthsStr)
	if months <= 0 {
		months = 6
	}

	results, err := c.Service.GetCashFlowHistory(tenantID, months)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, results)
}
