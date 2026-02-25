package models

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// GenerateULID is a helper to generate a new ULID string.
func GenerateULID() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id, _ := ulid.New(ulid.Timestamp(time.Now()), entropy)
	return id.String()
}

// CustomerCredit tracks credit limits and balances for customers.
type CustomerCredit struct {
	ID           string         `gorm:"primaryKey;type:char(26)" json:"id"`
	CustomerName string         `json:"customer_name"`
	CreditLimit  float64        `json:"credit_limit"`
	UsedCredit   float64        `json:"used_credit"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *CustomerCredit) BeforeCreate(tx *gorm.DB) error {
	m.ID = GenerateULID()
	return nil
}

// Invoice represents a billing document.
type Invoice struct {
	ID          string           `gorm:"primaryKey;type:char(26)" json:"id"`
	Number      string           `gorm:"uniqueIndex" json:"number"`
	InvoiceDate time.Time        `json:"invoice_date"`
	DueDate     time.Time        `json:"due_date"`
	TotalAmount float64          `json:"total_amount"`
	Status      string           `json:"status"` // e.g., Draft, Sent, Paid, Overdue
	Details     []InvoiceDetail  `gorm:"foreignKey:InvoiceID" json:"details,omitempty"`
	Payments    []InvoicePayment `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (m *Invoice) BeforeCreate(tx *gorm.DB) error {
	m.ID = GenerateULID()
	return nil
}

// InvoiceDetail represents line items within an invoice.
type InvoiceDetail struct {
	ID          string         `gorm:"primaryKey;type:char(26)" json:"id"`
	InvoiceID   string         `gorm:"type:char(26)" json:"invoice_id"`
	Description string         `json:"description"`
	Quantity    int            `json:"quantity"`
	UnitPrice   float64        `json:"unit_price"`
	Subtotal    float64        `json:"subtotal"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *InvoiceDetail) BeforeCreate(tx *gorm.DB) error {
	m.ID = GenerateULID()
	return nil
}

// InvoicePayment tracks payments made against an invoice.
type InvoicePayment struct {
	ID           string               `gorm:"primaryKey;type:char(26)" json:"id"`
	InvoiceID    string               `gorm:"type:char(26)" json:"invoice_id"`
	Amount       float64              `json:"amount"`
	PaymentDate  time.Time            `json:"payment_date"`
	Method       string               `json:"method"` // e.g., Cash, Bank Transfer, Credit Card
	Transactions []CompanyTransaction `gorm:"foreignKey:ReferenceID" json:"transactions,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	DeletedAt    gorm.DeletedAt       `gorm:"index" json:"-"`
}

func (m *InvoicePayment) BeforeCreate(tx *gorm.DB) error {
	m.ID = GenerateULID()
	return nil
}

// CompanyTransaction represents the general ledger/financial transactions.
type CompanyTransaction struct {
	ID            string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Date          time.Time      `json:"date"`
	Description   string         `json:"description"`
	Amount        float64        `json:"amount"`
	Type          string         `json:"type"` // e.g., Debit, Credit
	ReferenceID   string         `gorm:"type:char(26)" json:"reference_id"`
	ReferenceType string         `json:"reference_type"` // e.g., InvoicePayment, Requisition
	Requisitions  []Requisition  `gorm:"foreignKey:TransactionID" json:"requisitions,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *CompanyTransaction) BeforeCreate(tx *gorm.DB) error {
	m.ID = GenerateULID()
	return nil
}

// Requisition represents a request to buy goods or services (formerly Buy Request).
type Requisition struct {
	ID            string         `gorm:"primaryKey;type:char(26)" json:"id"`
	TransactionID string         `gorm:"type:char(26)" json:"transaction_id"`
	Description   string         `json:"description"`
	Amount        float64        `json:"amount"`
	Status        string         `json:"status"` // e.g., Pending, Approved, Completed
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Requisition) BeforeCreate(tx *gorm.DB) error {
	m.ID = GenerateULID()
	return nil
}

// CompanyBalance tracks various accounts (Bank, Cash, etc.)
type CompanyBalance struct {
	ID          string         `gorm:"primaryKey;type:char(26)" json:"id"`
	AccountName string         `json:"account_name"` // e.g., "Main Bank Account", "Petty Cash"
	Balance     float64        `json:"balance"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *CompanyBalance) BeforeCreate(tx *gorm.DB) error {
	m.ID = GenerateULID()
	return nil
}
