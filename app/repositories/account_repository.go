package repositories

import "gosvelte/app/models"

type AccountRepository struct {
	BaseRepository[models.Account]
}
