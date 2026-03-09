package repositories

import "gosvelte/app/models"

type TenantRepository struct {
	BaseRepository[models.Tenant]
}
