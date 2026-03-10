package controllers

import (
	"gosvelte/app/middleware"
	"gosvelte/app/services"
	"net/http"
	"time"
)

type AuthController struct {
	BaseController
	Service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{
		Service: service,
	}
}

// setSessionCookie is a helper to set the session_id cookie.
func (c *AuthController) setSessionCookie(w http.ResponseWriter, sessionID string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

// RegisterRequest defines the input for the registration endpoint.
type RegisterRequest struct {
	TenantName string `json:"tenant_name"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

// Register handles user registration and sets a session cookie.
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := c.Bind(r, &req); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" || req.TenantName == "" {
		c.Error(w, http.StatusBadRequest, "missing required fields")
		return
	}

	user, err := c.Service.Register(req.TenantName, req.Name, req.Email, req.Password)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create session automatically after registration
	session, err := c.Service.CreateSession(user.ID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	c.setSessionCookie(w, session.ID, session.ExpiresAt)
	c.JSON(w, http.StatusCreated, user)
}

// LoginRequest defines the input for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login authenticates a user and sets a session cookie.
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := c.Bind(r, &req); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := c.Service.Login(req.Email, req.Password)
	if err != nil {
		c.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	session, err := c.Service.CreateSession(user.ID)
	if err != nil {
		c.Error(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	c.setSessionCookie(w, session.ID, session.ExpiresAt)
	c.JSON(w, http.StatusOK, user)
}

// Logout revokes the current session and clears the cookie.
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		c.Service.Logout(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	c.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// GetProfile retrieves the currently logged-in user's information from the context.
func (c *AuthController) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		c.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := c.Service.GetContext(userID)
	if err != nil {
		c.Error(w, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(w, http.StatusOK, user)
}
