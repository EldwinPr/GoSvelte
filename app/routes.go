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
	mux.HandleFunc("POST /api/explorer", mw.Auth(mw.RequireClearance(20, explorerController.CreateRecord)))
	mux.HandleFunc("PUT /api/explorer", mw.Auth(mw.RequireClearance(20, explorerController.UpdateRecord)))
	mux.HandleFunc("DELETE /api/explorer", mw.Auth(mw.RequireClearance(20, explorerController.DeleteRecord)))

	// Dashboard
	mux.HandleFunc("GET /api/dashboard/stats", mw.Auth(dashboardController.Stats))

	// Users (Developer Only)
	mux.HandleFunc("GET /api/users", mw.Auth(mw.RequireClearance(20, userController.Index)))

	// Invoices (Finance 0+)
	mux.HandleFunc("GET /api/invoices", mw.Auth(invoiceController.Index))
	mux.HandleFunc("GET /api/invoices/show", mw.Auth(invoiceController.Show))
	mux.HandleFunc("GET /api/payments", mw.Auth(invoiceController.Payments))
	mux.HandleFunc("GET /api/payments/summary", mw.Auth(invoiceController.PaymentSummary))
	mux.HandleFunc("POST /api/invoices", mw.Auth(invoiceController.Create))
	mux.HandleFunc("PUT /api/invoices", mw.Auth(invoiceController.Update))
	mux.HandleFunc("POST /api/invoices/pay", mw.Auth(invoiceController.Pay))

	// Requisitions
	mux.HandleFunc("GET /api/requisitions", mw.Auth(requisitionController.Index))
	mux.HandleFunc("GET /api/requisitions/show", mw.Auth(requisitionController.Show))
	mux.HandleFunc("POST /api/requisitions", mw.Auth(requisitionController.Create))
	// Approval (Manager 10+)
	mux.HandleFunc("POST /api/requisitions/approve", mw.Auth(mw.RequireClearance(10, requisitionController.Approve)))
	mux.HandleFunc("POST /api/requisitions/reject", mw.Auth(mw.RequireClearance(10, requisitionController.Reject)))
	// Disbursement (Finance 0+)
	mux.HandleFunc("POST /api/requisitions/give", mw.Auth(mw.RequireClearance(0, requisitionController.MarkGiven)))

	// Transactions (Finance 0+)
	mux.HandleFunc("GET /api/transactions", mw.Auth(transactionController.Index))

	// Balance (Manager 10+)
	mux.HandleFunc("GET /api/balances", mw.Auth(mw.RequireClearance(0, balanceController.Index)))

	// --- 3. Static & SPA Routes ---
	fileServer := http.FileServer(http.Dir("./static"))
	
	// Cache-Control middleware for hashed static assets (JS/CSS)
	staticWithCache := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable") 
		fileServer.ServeHTTP(w, r)
	})

	mux.Handle("/static/", http.StripPrefix("/static/", staticWithCache))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "./static" + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// index.html should not be cached long to ensure users get updates
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			http.ServeFile(w, r, "./static/index.html")
			return
		}
		// Root assets (icons, etc)
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fileServer.ServeHTTP(w, r)
	})

	return mux
}
