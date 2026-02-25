package app

import (
	"net/http"
	"os"

	"gosvelte/app/controllers"

	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mw := NewMiddleware(db)

	// --- 1. Controllers ---
	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	invoiceController := controllers.NewInvoiceController(db)
	requisitionController := controllers.NewRequisitionController(db)
	transactionController := controllers.NewTransactionController(db)
	balanceController := controllers.NewBalanceController(db)
	dashboardController := controllers.NewDashboardController(db)
	explorerController := controllers.NewExplorerController(db)

	// --- 2. API Routes ---
	
	// Auth
	mux.HandleFunc("POST /api/login", authController.Login)
	mux.HandleFunc("POST /api/logout", authController.Logout)
	mux.HandleFunc("GET /api/me", mw.Auth(authController.Me))
	mux.HandleFunc("POST /api/register", mw.Auth(mw.RequireClearance(20, authController.Register)))

	// Explorer
	mux.HandleFunc("GET /api/explorer", mw.Auth(mw.RequireClearance(20, explorerController.GetTableData)))

	// Dashboard
	mux.HandleFunc("GET /api/dashboard/stats", mw.Auth(dashboardController.Stats))

	// Users (Developer Only)
	mux.HandleFunc("GET /api/users", mw.Auth(mw.RequireClearance(20, userController.Index)))

	// Invoices (Finance 0+)
	mux.HandleFunc("GET /api/invoices", mw.Auth(invoiceController.Index))
	mux.HandleFunc("POST /api/invoices", mw.Auth(invoiceController.Create))
	mux.HandleFunc("POST /api/invoices/pay", mw.Auth(invoiceController.Pay))

	// Requisitions
	mux.HandleFunc("GET /api/requisitions", mw.Auth(requisitionController.Index))
	mux.HandleFunc("POST /api/requisitions", mw.Auth(requisitionController.Create))
	// Approval (Manager 10+)
	mux.HandleFunc("POST /api/requisitions/approve", mw.Auth(mw.RequireClearance(10, requisitionController.Approve)))

	// Transactions (Finance 0+)
	mux.HandleFunc("GET /api/transactions", mw.Auth(transactionController.Index))

	// Balance (Manager 10+)
	mux.HandleFunc("GET /api/balances", mw.Auth(mw.RequireClearance(10, balanceController.Index)))

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
