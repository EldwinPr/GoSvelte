package main

import (
	"log"
	"net/http"
	"os"

	"gosvelte/app"
	"gosvelte/app/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// Database Initialization
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Fallback for local development if not provided
		dsn = "host=localhost user=user password=pass dbname=erp port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-Migrate
	db.AutoMigrate(&models.User{})

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
