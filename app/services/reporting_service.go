package services

import (
	"gosvelte/app/models"
	"time"

	"gorm.io/gorm"
)

type ReportingService struct {
	BaseService
}

func NewReportingService(db *gorm.DB) *ReportingService {
	return &ReportingService{
		BaseService: BaseService{DB: db},
	}
}

// SpendingBreakdown represents total spending per category.
type SpendingBreakdown struct {
	CategoryName string  `json:"category_name"`
	Color        string  `json:"color"`
	Total        float64 `json:"total"`
}

// GetSpendingBreakdown calculates total expenses per category for a date range.
func (s *ReportingService) GetSpendingBreakdown(tenantID string, start, end time.Time) ([]SpendingBreakdown, error) {
	var results []SpendingBreakdown
	err := s.DB.Model(&models.Transaction{}).
		Select("categories.name as category_name, categories.color, sum(transactions.amount) as total").
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.tenant_id = ? AND transactions.type = ? AND transactions.date >= ? AND transactions.date <= ?", 
			tenantID, "Expense", start, end).
		Group("categories.name, categories.color").
		Order("total DESC").
		Scan(&results).Error
	
	return results, err
}

// MonthlyCashFlow represents income vs expense for a specific month.
type MonthlyCashFlow struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

// GetCashFlowHistory retrieves income vs expense data for the last N months.
func (s *ReportingService) GetCashFlowHistory(tenantID string, months int) ([]MonthlyCashFlow, error) {
	var results []MonthlyCashFlow
	
	// We'll generate a query that groups by month
	// Note: PostgreSQL specific date formatting. 
	err := s.DB.Model(&models.Transaction{}).
		Select("to_char(date, 'YYYY-MM') as month, " +
			"sum(case when type = 'Income' then amount else 0 end) as income, " +
			"sum(case when type = 'Expense' then amount else 0 end) as expense").
		Where("tenant_id = ? AND date >= ?", tenantID, time.Now().AddDate(0, -months, 0)).
		Group("month").
		Order("month ASC").
		Scan(&results).Error

	return results, err
}

// DashboardSnapshot is a composite model for the main dashboard view.
type DashboardSnapshot struct {
	NetWorth          float64             `json:"net_worth"`
	TotalLent         float64             `json:"total_lent"`
	TotalBorrowed     float64             `json:"total_borrowed"`
	RecentActivity    []models.Transaction `json:"recent_activity"`
	SpendingBreakdown []SpendingBreakdown `json:"spending_breakdown"`
	MonthlyCashFlow   []MonthlyCashFlow   `json:"monthly_cash_flow"`
}

// GetDashboardSnapshot aggregates data from multiple sources for a single "loading" state.
func (s *ReportingService) GetDashboardSnapshot(tenantID string) (*DashboardSnapshot, error) {
	snapshot := &DashboardSnapshot{}

	// 1. Calculate Net Worth
	s.DB.Model(&models.Account{}).
		Where("tenant_id = ?", tenantID).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&snapshot.NetWorth)

	// 2. Calculate Debt Summary
	s.DB.Model(&models.Debt{}).
		Select("sum(case when type = 'Lent' then remaining else 0 end) as total_lent, " +
			"sum(case when type = 'Borrowed' then remaining else 0 end) as total_borrowed").
		Where("tenant_id = ?", tenantID).
		Row().Scan(&snapshot.TotalLent, &snapshot.TotalBorrowed)

	// 3. Fetch Recent Activity (last 5)
	s.DB.Preload("Account").Preload("Category").
		Where("tenant_id = ?", tenantID).
		Order("date DESC").Limit(5).
		Find(&snapshot.RecentActivity)

	// 4. Get Current Month Spending Breakdown
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	snapshot.SpendingBreakdown, _ = s.GetSpendingBreakdown(tenantID, startOfMonth, now)

	// 5. Get 6-month Cash Flow
	snapshot.MonthlyCashFlow, _ = s.GetCashFlowHistory(tenantID, 6)

	return snapshot, nil
}
