package controllers

import (
	"net/http"

	"gorm.io/gorm"
)

type ExplorerController struct {
	BaseController
	DB *gorm.DB
}

func NewExplorerController(db *gorm.DB) *ExplorerController {
	return &ExplorerController{DB: db}
}

func (c *ExplorerController) GetTableData(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	if table == "" {
		c.Error(w, http.StatusBadRequest, "Table name is required")
		return
	}

	// For security, only allow specific tables even for developers
	allowedTables := map[string]bool{
		"users":                true,
		"customer_credits":      true,
		"invoices":             true,
		"invoice_details":      true,
		"invoice_payments":     true,
		"company_transactions": true,
		"requisitions":         true,
		"company_balances":      true,
	}

	if !allowedTables[table] {
		c.Error(w, http.StatusForbidden, "Table not allowed or does not exist")
		return
	}

	results := []map[string]interface{}{}
	if err := c.DB.Table(table).Find(&results).Error; err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, results)
}
