package services

import (
	"gosvelte/app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Auto-Migrate all models
	db.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.Account{},
		&models.Category{},
		&models.Transaction{},
		&models.Transfer{},
		&models.Debt{},
		&models.DebtInstallment{},
		&models.RecurringPayment{},
		&models.RecurringInstance{},
		&models.Wishlist{},
		&models.Budget{},
		&models.UserSettings{},
		&models.Session{},
	)

	return db
}
