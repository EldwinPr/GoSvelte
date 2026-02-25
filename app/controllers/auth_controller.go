package controllers

import (
	"encoding/json"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type AuthController struct {
	BaseController
	Service *services.AuthService
}

func NewAuthController(db *gorm.DB) *AuthController {
	repo := &repositories.UserRepository{BaseRepository: repositories.BaseRepository[models.User]{DB: db}}
	service := &services.AuthService{BaseService: services.BaseService{DB: db}, Repo: repo}
	return &AuthController{Service: service}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Service.Register(&user); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(w, http.StatusCreated, user)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := c.Service.Login(input.Email, input.Password)
	if err != nil {
		c.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set a session cookie (In prototype, we use User ID as token)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user.ID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(w, http.StatusOK, user)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	c.JSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
}

func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	// Import "gosvelte/app" is needed or use the key correctly
	// Since we are in controllers, we can't import app (circular dependency)
	// So we should move the key to a shared package or just use string in both places for the prototype.
	// Let's use string in both places for simplicity in this prototype.
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		c.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}
	c.JSON(w, http.StatusOK, user)
}
