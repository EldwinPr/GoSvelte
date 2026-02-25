package repositories

import "gosvelte/app/models"

type InvoicePaymentRepository struct {
	BaseRepository[models.InvoicePayment]
}
