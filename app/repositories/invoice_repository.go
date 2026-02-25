package repositories

import "gosvelte/app/models"

type InvoiceRepository struct {
	BaseRepository[models.Invoice]
}
