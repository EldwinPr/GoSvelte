package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// BaseController provides standardized response helpers for HTTP handlers.
type BaseController struct{}

// JSON responds with a standardized JSON structure.
func (c *BaseController) JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error responds with a standardized error message.
func (c *BaseController) Error(w http.ResponseWriter, status int, message string) {
	c.JSON(w, status, map[string]string{"error": message})
}

// QueryOptions holds standard query parameters for filtering, sorting, and pagination.
type QueryOptions struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
	Search   string
}

// GetOrder returns the order string for GORM (e.g. "name asc")
func (o QueryOptions) GetOrder(defaultSort string) string {
	if o.Sort == "" {
		return defaultSort
	}
	order := o.Order
	if order == "" {
		order = "asc"
	}
	return o.Sort + " " + order
}

// ParseQuery parses standard query parameters from the request.
func (c *BaseController) ParseQuery(r *http.Request) QueryOptions {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	return QueryOptions{
		Page:     page,
		PageSize: pageSize,
		Sort:     r.URL.Query().Get("sort"),
		Order:    r.URL.Query().Get("order"),
		Search:   r.URL.Query().Get("search"),
	}
}
