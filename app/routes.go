package app

import (
	"net/http"
	"os"

	"gosvelte/app/controllers"

	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// --- 1. Controllers (Internal DI is handled in constructors) ---
	userController := controllers.NewUserController(db)

	// --- 2. API Routes ---
	mux.HandleFunc("GET /api/users", userController.Index)

	// --- 3. Static & SPA Routes ---
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "./static" + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, "./static/index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	return mux
}
