package services

import (
	"encoding/json"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"time"

	"gorm.io/gorm"
)

type AuditService struct {
	BaseService
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{
		BaseService: BaseService{DB: db},
	}
}

// Log records a change in the system.
func (s *AuditService) Log(tenantID, userID, action, entityType, entityID string, oldVal, newVal interface{}) error {
	var oldJSON, newJSON string
	
	if oldVal != nil {
		b, _ := json.Marshal(oldVal)
		oldJSON = string(b)
	}
	
	if newVal != nil {
		b, _ := json.Marshal(newVal)
		newJSON = string(b)
	}

	log := &models.AuditLog{
		ID:         repositories.GenerateULID(),
		TenantID:   tenantID,
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValue:   oldJSON,
		NewValue:   newJSON,
		CreatedAt:  time.Now(),
	}

	return s.DB.Create(log).Error
}
