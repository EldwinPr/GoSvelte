package repositories

import "gosvelte/app/models"

type CompanyTransactionRepository struct {
	BaseRepository[models.CompanyTransaction]
}
