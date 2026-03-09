package repositories

import "gosvelte/app/models"

type CategoryRepository struct {
	BaseRepository[models.Category]
}
