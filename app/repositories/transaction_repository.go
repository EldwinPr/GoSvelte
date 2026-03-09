package repositories

import "gosvelte/app/models"

type TransactionRepository struct {
	BaseRepository[models.Transaction]
}
