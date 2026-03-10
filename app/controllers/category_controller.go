package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/models"
	"gosvelte/app/services"
	"net/http"
)

type CategoryController struct {
	BaseController
	Service *services.CategoryService
}

func NewCategoryController(service *services.CategoryService) *CategoryController {
	return &CategoryController{
		Service: service,
	}
}

// List retrieves all categories for the current tenant.
func (c *CategoryController) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	categories, err := c.Service.GetCategories(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, categories)
}

// Create handles the creation of a new category.
func (c *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var category models.Category
	if err := c.Bind(r, &category); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.CreateCategory(tenantID, &category); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, category)
}

// Update handles category modifications.
func (c *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var category models.Category
	if err := c.Bind(r, &category); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if category.ID == "" {
		c.Error(w, http.StatusBadRequest, "category id is required")
		return
	}

	if err := c.Service.UpdateCategory(tenantID, &category); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, category)
}

// Delete handles category removal.
func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id")

	if err := c.Service.DeleteCategory(tenantID, id); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusNoContent, nil)
}
