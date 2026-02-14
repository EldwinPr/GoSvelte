package controllers

import (
	"encoding/json"
	"net/http"
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
