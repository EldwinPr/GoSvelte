package main

import (
	"log"
	"net/http"
	"os"

	"gosvelte/app"
	"gosvelte/app/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	// "gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// Database Initialization
	// PostgreSQL settings (commented out)
	/*
		dsn := os.Getenv("DB_DSN")
		if dsn == "" {
			// Fallback for local development if not provided
			dsn = "host=localhost user=user password=pass dbname=erp port=5432 sslmode=disable"
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	*/

	// SQLite settings
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH environment variable is required but not set")
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-Migrate
	db.AutoMigrate(
		&models.User{},
		&models.CustomerCredit{},
		&models.Invoice{},
		&models.InvoiceDetail{},
		&models.InvoicePayment{},
		&models.CompanyTransaction{},
		&models.Requisition{},
		&models.CompanyBalance{},
	)

	// Seed Initial Data (Prototype only)
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		admin := &models.User{
			Name:      "Admin Developer",
			Email:     "admin@sentral.com",
			Password:  string(hashedPassword),
			Clearance: 20,
		}
		db.Create(admin)
		log.Println("Created default admin: admin@sentral.com / password")
	}

	db.Model(&models.CompanyBalance{}).Count(&count)
	if count == 0 {
		mainBank := &models.CompanyBalance{
			AccountName: "OCBC",
			Balance:     1000000.0,
		}
		db.Create(mainBank)
		log.Println("Created initial company balance: 1,000,000.0")
	}

	// Register Routes and Start Server
	mux := app.RegisterRoutes(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
