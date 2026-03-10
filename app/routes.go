package app

import (
	"net/http"
	"os"

	"gosvelte/app/controllers"
	"gosvelte/app/middleware"
	"gosvelte/app/services"

	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()

	// --- 1. Instantiate Services ---
	authService := services.NewAuthService(db)
	accountService := services.NewAccountService(db)
	transactionService := services.NewTransactionService(db)
	categoryService := services.NewCategoryService(db)
	transferService := services.NewTransferService(db)
	recurringService := services.NewRecurringPaymentService(db, transactionService)
	debtService := services.NewDebtService(db)
	wishlistService := services.NewWishlistService(db, accountService)
	budgetService := services.NewBudgetService(db)
	reportingService := services.NewReportingService(db)

	// --- 2. Instantiate Controllers with Injected Services ---
	authController := controllers.NewAuthController(authService)
	accountController := controllers.NewAccountController(accountService)
	transactionController := controllers.NewTransactionController(transactionService)
	categoryController := controllers.NewCategoryController(categoryService)
	transferController := controllers.NewTransferController(transferService)
	recurringController := controllers.NewRecurringPaymentController(recurringService)
	debtController := controllers.NewDebtController(debtService)
	wishlistController := controllers.NewWishlistController(wishlistService)
	budgetController := controllers.NewBudgetController(budgetService)
	reportingController := controllers.NewReportingController(reportingService)

	// --- 3. Public API Routes ---
	mux.HandleFunc("POST /api/auth/register", authController.Register)
	mux.HandleFunc("POST /api/auth/login", authController.Login)
	mux.HandleFunc("POST /api/auth/logout", authController.Logout)

	// --- 4. Protected API Routes ---
	authMid := middleware.AuthMiddleware(authService)
	
	// Auth Profile
	mux.Handle("GET /api/auth/profile", authMid(http.HandlerFunc(authController.GetProfile)))

	// Account Routes
	mux.Handle("GET /api/accounts", authMid(http.HandlerFunc(accountController.List)))
	mux.Handle("POST /api/accounts", authMid(http.HandlerFunc(accountController.Create)))
	mux.Handle("GET /api/accounts/detail", authMid(http.HandlerFunc(accountController.Get)))
	mux.Handle("PUT /api/accounts", authMid(http.HandlerFunc(accountController.Update)))
	mux.Handle("DELETE /api/accounts", authMid(http.HandlerFunc(accountController.Delete)))
	mux.Handle("GET /api/accounts/dashboard", authMid(http.HandlerFunc(accountController.DashboardSummary)))

	// Transaction Routes
	mux.Handle("GET /api/transactions", authMid(http.HandlerFunc(transactionController.List)))
	mux.Handle("POST /api/transactions", authMid(http.HandlerFunc(transactionController.Create)))
	mux.Handle("PUT /api/transactions", authMid(http.HandlerFunc(transactionController.Update)))
	mux.Handle("DELETE /api/transactions", authMid(http.HandlerFunc(transactionController.Delete)))

	// Category Routes
	mux.Handle("GET /api/categories", authMid(http.HandlerFunc(categoryController.List)))
	mux.Handle("POST /api/categories", authMid(http.HandlerFunc(categoryController.Create)))
	mux.Handle("PUT /api/categories", authMid(http.HandlerFunc(categoryController.Update)))
	mux.Handle("DELETE /api/categories", authMid(http.HandlerFunc(categoryController.Delete)))

	// Transfer Routes
	mux.Handle("GET /api/transfers", authMid(http.HandlerFunc(transferController.List)))
	mux.Handle("POST /api/transfers", authMid(http.HandlerFunc(transferController.Create)))
	mux.Handle("DELETE /api/transfers", authMid(http.HandlerFunc(transferController.Delete)))

	// Recurring Payment Routes
	mux.Handle("GET /api/recurring", authMid(http.HandlerFunc(recurringController.List)))
	mux.Handle("POST /api/recurring", authMid(http.HandlerFunc(recurringController.Create)))
	mux.Handle("POST /api/recurring/toggle", authMid(http.HandlerFunc(recurringController.Toggle)))
	mux.Handle("POST /api/recurring/check", authMid(http.HandlerFunc(recurringController.ProcessCheck)))
	mux.Handle("GET /api/recurring/pending", authMid(http.HandlerFunc(recurringController.ListPending)))
	mux.Handle("POST /api/recurring/confirm", authMid(http.HandlerFunc(recurringController.Confirm)))
	mux.Handle("POST /api/recurring/skip", authMid(http.HandlerFunc(recurringController.Skip)))

	// Debt Routes
	mux.Handle("GET /api/debts", authMid(http.HandlerFunc(debtController.List)))
	mux.Handle("POST /api/debts", authMid(http.HandlerFunc(debtController.Create)))
	mux.Handle("POST /api/debts/installment", authMid(http.HandlerFunc(debtController.AddInstallment)))
	mux.Handle("GET /api/debts/summary", authMid(http.HandlerFunc(debtController.Summary)))

	// Wishlist Routes
	mux.Handle("GET /api/wishlist", authMid(http.HandlerFunc(wishlistController.List)))
	mux.Handle("POST /api/wishlist", authMid(http.HandlerFunc(wishlistController.Create)))
	mux.Handle("POST /api/wishlist/progress", authMid(http.HandlerFunc(wishlistController.UpdateProgress)))
	mux.Handle("POST /api/wishlist/purchase", authMid(http.HandlerFunc(wishlistController.Purchase)))
	mux.Handle("DELETE /api/wishlist", authMid(http.HandlerFunc(wishlistController.Delete)))

	// Budget Routes
	mux.Handle("GET /api/budgets", authMid(http.HandlerFunc(budgetController.List)))
	mux.Handle("POST /api/budgets", authMid(http.HandlerFunc(budgetController.Create)))
	mux.Handle("PUT /api/budgets", authMid(http.HandlerFunc(budgetController.Update)))
	mux.Handle("DELETE /api/budgets", authMid(http.HandlerFunc(budgetController.Delete)))
	mux.Handle("GET /api/budgets/performance", authMid(http.HandlerFunc(budgetController.Performance)))

	// Reporting & Dashboard Routes
	mux.Handle("GET /api/reports/dashboard", authMid(http.HandlerFunc(reportingController.GetDashboard)))
	mux.Handle("GET /api/reports/spending", authMid(http.HandlerFunc(reportingController.GetSpendingBreakdown)))
	mux.Handle("GET /api/reports/cashflow", authMid(http.HandlerFunc(reportingController.GetCashFlowHistory)))

	// --- 5. Static & SPA Routes ---
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

	// Wrap the entire mux with global middleware
	handler := middleware.RecoveryMiddleware(mux)
	handler = middleware.LoggerMiddleware(handler)

	return handler
}
