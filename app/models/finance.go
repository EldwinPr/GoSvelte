package models

import (
	"gorm.io/gorm"
	"time"
)

type Tenant struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	TenantID  string         `gorm:"index" json:"tenant_id"`
	Tenant    Tenant         `json:"tenant,omitempty"`
	Name      string         `json:"name"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Account struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	TenantID  string         `gorm:"index" json:"tenant_id"`
	Name      string         `json:"name"`
	Type      string         `json:"type"` // e.g., Bank, Cash, Credit Card
	Balance   float64        `json:"balance" gorm:"default:0"`
	Currency  string         `json:"currency" gorm:"default:'USD'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Category struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	TenantID  string         `gorm:"index" json:"tenant_id"`
	Name      string         `json:"name"`
	Icon      string         `json:"icon"`
	Color     string         `json:"color"`
	Type      string         `json:"type"` // e.g., Income, Expense
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transfer struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	TenantID      string         `gorm:"index" json:"tenant_id"`
	FromAccountID string         `gorm:"index" json:"from_account_id"`
	ToAccountID   string         `gorm:"index" json:"to_account_id"`
	FromAccount   Account        `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
	ToAccount     Account        `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
	Transactions  []Transaction  `json:"transactions,omitempty"`
	Amount        float64        `json:"amount"`
	Description   string         `json:"description"`
	Date          time.Time      `json:"date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type RecurringPayment struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"index" json:"tenant_id"`
	AccountID   *string        `gorm:"index" json:"account_id"`
	CategoryID  string         `gorm:"index" json:"category_id"`
	Account     *Account       `json:"account,omitempty"`
	Category    Category       `json:"category,omitempty"`
	Amount      float64        `json:"amount"`
	Description string         `json:"description"`
	Frequency   string         `json:"frequency"` // e.g., Monthly, Weekly
	StartDate   time.Time      `json:"start_date"`
	NextDate    time.Time      `json:"next_date"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Debt struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"index" json:"tenant_id"`
	Name        string         `json:"name"`
	TotalAmount float64        `json:"total_amount"`
	Remaining   float64        `json:"remaining"`
	DueDate     time.Time      `json:"due_date"`
	Type        string         `json:"type"` // e.g., Lent, Borrowed
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type DebtInstallment struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	TenantID  string         `gorm:"index" json:"tenant_id"`
	DebtID    string         `gorm:"index" json:"debt_id"`
	Debt      Debt           `json:"debt,omitempty"`
	Amount    float64        `json:"amount"`
	PaidDate  time.Time      `json:"paid_date"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	ID                 string           `gorm:"primaryKey" json:"id"`
	TenantID           string           `gorm:"index" json:"tenant_id"`
	AccountID          string           `gorm:"index" json:"account_id"`
	CategoryID         string           `gorm:"index" json:"category_id"`
	TransferID         *string          `gorm:"index" json:"transfer_id,omitempty"`
	RecurringPaymentID *string          `gorm:"index" json:"recurring_payment_id,omitempty"`
	DebtID             *string          `gorm:"index" json:"debt_id,omitempty"`
	DebtInstallmentID  *string          `gorm:"index" json:"debt_installment_id,omitempty"`
	
	Account            Account          `json:"account,omitempty"`
	Category           Category         `json:"category,omitempty"`
	Transfer           *Transfer        `json:"transfer,omitempty"`
	RecurringPayment   *RecurringPayment `json:"recurring_payment,omitempty"`
	Debt               *Debt            `json:"debt,omitempty"`
	DebtInstallment    *DebtInstallment  `json:"debt_installment,omitempty"`

	Amount             float64          `json:"amount"`
	Description        string           `json:"description"`
	Date               time.Time        `json:"date"`
	Type               string           `json:"type"` // Income, Expense
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	DeletedAt          gorm.DeletedAt   `gorm:"index" json:"-"`
}

type Wishlist struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"index" json:"tenant_id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	TargetAmount float64        `json:"target_amount"`
	CurrentSaved float64        `json:"current_saved" gorm:"default:0"`
	Need         int            `json:"need" gorm:"default:1"`         // 1-6 scale
	Want         int            `json:"want" gorm:"default:1"`         // 1-6 scale
	Productivity int            `json:"productivity" gorm:"default:1"` // 1-6 scale
	URL          string         `json:"url"`
	Status       string         `json:"status" gorm:"default:'Active'"` // Active, Purchased, Archived
	DesiredDate  *time.Time     `json:"desired_date"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Virtual Fields (Calculated at runtime, not saved in DB)
	FinancialImpact float64 `gorm:"-" json:"financial_impact"`
	Score           float64 `gorm:"-" json:"score"`
}

type Budget struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	TenantID   string         `gorm:"index" json:"tenant_id"`
	CategoryID *string        `gorm:"index" json:"category_id"` // Null = Global Budget
	Category   *Category      `json:"category,omitempty"`
	Amount     float64        `json:"amount"`
	Period     string         `json:"period"` // Monthly, Weekly, Yearly, One-Time
	StartDate  time.Time      `json:"start_date"`
	EndDate    *time.Time     `json:"end_date"`
	IsStrict   bool           `gorm:"default:false" json:"is_strict"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type SavingsGoal struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	TenantID      string         `gorm:"index" json:"tenant_id"`
	Name          string         `json:"name"`
	TargetAmount  float64        `json:"target_amount"`
	CurrentAmount float64        `json:"current_amount" gorm:"default:0"`
	AccountID     *string        `gorm:"index" json:"account_id,omitempty"` // Optional dedicated account
	Account       *Account       `json:"account,omitempty"`
	TargetDate    *time.Time     `json:"target_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
