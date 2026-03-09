package repositories

import "gosvelte/app/models"

type UserRepository struct {
	BaseRepository[models.User]
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}
