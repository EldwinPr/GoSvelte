package services

import "gorm.io/gorm"

// BaseService provides shared functionality for business logic layers.
type BaseService struct {
	DB *gorm.DB
}

// Transaction wraps a function in a database transaction.
func (s *BaseService) Transaction(fn func(tx *gorm.DB) error) error {
	return s.DB.Transaction(fn)
}
