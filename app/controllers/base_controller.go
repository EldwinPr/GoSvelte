package controllers

import (
	"encoding/json"
	"gosvelte/app/repositories"
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

// Bind parses the request body into the provided interface.
func (c *BaseController) Bind(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// GetPaginationParams parses pagination and sorting parameters from the request query.
func (c *BaseController) GetPaginationParams(r *http.Request) *repositories.Pagination {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	page, _ := strconv.Atoi(query.Get("page"))
	sort := query.Get("sort")

	return &repositories.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

// PaginatedResponse responds with a standardized paginated JSON structure.
func (c *BaseController) PaginatedResponse(w http.ResponseWriter, p *repositories.Pagination) {
	c.JSON(w, http.StatusOK, p)
}
