package controllers

import (
	"encoding/json"
	"gosvelte/app/models"
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

var allowedTables = map[string]bool{
	"users":                true,
	"customer_credits":      true,
	"invoices":             true,
	"invoice_details":      true,
	"invoice_payments":     true,
	"company_transactions": true,
	"requisitions":         true,
	"company_balances":      true,
}

func (c *ExplorerController) GetTableData(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	if !allowedTables[table] {
		c.Error(w, http.StatusForbidden, "Table not allowed or does not exist")
		return
	}

	opts := c.ParseQuery(r)
	
	var totalCount int64
	query := c.DB.Table(table).Unscoped() // Include soft-deleted rows
	
	if opts.Search != "" {
		// Generic search attempt (only works if 'name' or 'description' columns exist)
		// For a real explorer, we'd need to query column names first
		query = query.Where("id LIKE ?", "%"+opts.Search+"%")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	results := []map[string]interface{}{}
	offset := (opts.Page - 1) * opts.PageSize
	
	orderStr := opts.GetOrder("id desc")
	
	if err := query.Order(orderStr).Limit(opts.PageSize).Offset(offset).Find(&results).Error; err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]interface{}{
		"items":       results,
		"total_count": totalCount,
		"page":        opts.Page,
		"page_size":   opts.PageSize,
	})
}

func (c *ExplorerController) CreateRecord(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	if !allowedTables[table] {
		c.Error(w, http.StatusForbidden, "Table not allowed")
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Manually generate ULID since map-based creation bypasses model hooks
	data["id"] = models.GenerateULID()
	
	// Remove other system fields to let DB/GORM handle them
	delete(data, "created_at")
	delete(data, "updated_at")
	delete(data, "deleted_at")

	if err := c.DB.Table(table).Create(&data).Error; err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, data)
}

func (c *ExplorerController) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	if !allowedTables[table] {
		c.Error(w, http.StatusForbidden, "Table not allowed")
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	id, ok := data["id"]
	if !ok {
		c.Error(w, http.StatusBadRequest, "ID is required for updates")
		return
	}

	// Remove system fields from the update map to avoid GORM errors
	delete(data, "id")
	delete(data, "created_at")
	delete(data, "updated_at")
	delete(data, "deleted_at")

	if err := c.DB.Table(table).Where("id = ?", id).Updates(data).Error; err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]string{"message": "Record updated successfully"})
}

func (c *ExplorerController) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	id := r.URL.Query().Get("id")
	if !allowedTables[table] || id == "" {
		c.Error(w, http.StatusBadRequest, "Table name and ID are required")
		return
	}

	// Perform soft delete by setting deleted_at
	// Since we are using generic Table(), we have to do it manually if we don't have the model
	if err := c.DB.Table(table).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]string{"message": "Record soft-deleted"})
}
