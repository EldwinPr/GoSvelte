package repositories

import "gosvelte/app/models"

type UserRepository struct {
	BaseRepository[models.User]
}

// Add user-specific queries here
func (r *UserRepository) GetActiveUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.Where("active = ?", true).Find(&users).Error
	return users, err
}
