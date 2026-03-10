package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/models"
	"gosvelte/app/services"
	"net/http"
)

type WishlistController struct {
	BaseController
	Service *services.WishlistService
}

func NewWishlistController(service *services.WishlistService) *WishlistController {
	return &WishlistController{
		Service: service,
	}
}

// List retrieves the ranked wishlist for the current tenant.
func (c *WishlistController) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	items, err := c.Service.GetRankedWishlist(tenantID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusOK, items)
}

// Create handles the addition of a new wishlist item.
func (c *WishlistController) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var item models.Wishlist
	if err := c.Bind(r, &item); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.CreateItem(tenantID, &item); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, item)
}

// UpdateProgressRequest defines the input for updating current savings for an item.
type UpdateProgressRequest struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

// UpdateProgress handles manual adjustment of saved funds for a wishlist item.
func (c *WishlistController) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var req UpdateProgressRequest
	if err := c.Bind(r, &req); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.UpdateProgress(tenantID, req.ID, req.Amount); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]string{"message": "progress updated"})
}

// PurchaseRequest defines the input for marking an item as purchased.
type PurchaseRequest struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
}

// Purchase handles the conversion of a wishlist item into a real transaction.
func (c *WishlistController) Purchase(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	var req PurchaseRequest
	if err := c.Bind(r, &req); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.MarkPurchased(tenantID, req.ID, req.AccountID); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, map[string]string{"message": "item marked as purchased and transaction recorded"})
}

// Delete handles the removal of a wishlist item.
func (c *WishlistController) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	id := r.URL.Query().Get("id")

	if err := c.Service.DeleteItem(tenantID, id); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusNoContent, nil)
}
