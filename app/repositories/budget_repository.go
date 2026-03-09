package repositories

import "gosvelte/app/models"

type BudgetRepository struct {
	BaseRepository[models.Budget]
}
