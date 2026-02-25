package controllers

import (
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"gosvelte/app/services"
	"net/http"

	"gorm.io/gorm"
)

type UserController struct {
	BaseController
	Service *services.UserService
}

func NewUserController(db *gorm.DB) *UserController {
	repo := &repositories.UserRepository{
		BaseRepository: repositories.BaseRepository[models.User]{DB: db},
	}
	service := &services.UserService{
		BaseService: services.BaseService{DB: db},
		Repo:        repo,
	}
	return &UserController{
		Service: service,
	}
}

func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	users, err := c.Service.GetAllUsers()
	if err != nil {
		c.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if users == nil {
		users = []models.User{}
	}
	c.JSON(w, http.StatusOK, users)
}
