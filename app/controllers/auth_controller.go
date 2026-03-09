package controllers

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"

	"gorm.io/gorm"
)

type AuthController struct {
	BaseController
	Service *services.AuthService
}

func NewAuthController(db *gorm.DB) *AuthController {
	repo := &repositories.UserRepository{
		BaseRepository: repositories.BaseRepository[models.User]{DB: db},
	}
	service := &services.AuthService{
		BaseService: services.BaseService{DB: db},
		Repo:        repo,
	}
	return &AuthController{
		Service: service,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := c.Bind(r, &user); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.Service.Register(&user); err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(w, http.StatusCreated, user)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(r, &credentials); err != nil {
		c.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := c.Service.Login(credentials.Email, credentials.Password)
	if err != nil {
		c.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(w, http.StatusOK, user)
}
