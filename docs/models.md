# Models Layer

The models layer defines the database schema using GORM structs. Each model represents a table in the PostgreSQL database.

## Overview

All models support:
- **Soft deletes** via `gorm.DeletedAt` (records are never permanently deleted)
- **Timestamps** (`CreatedAt`, `UpdatedAt`) automatically managed by GORM
- **Multi-tenancy** via `TenantID` on most models
- **JSON serialization** with proper field tags

---

## Core Models

### Tenant
The root entity for multi-tenant isolation.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key (ULID) |
| `Name` | string | Tenant name |
| `CreatedAt` | time.Time | Creation timestamp |
| `UpdatedAt` | time.Time | Last update timestamp |
| `DeletedAt` | gorm.DeletedAt | Soft delete timestamp |

---

### User
Represents a user within a tenant.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `Name` | string | User's display name |
| `Email` | string | Unique email address |
| `Password` | string | Hashed password (excluded from JSON via `json:"-"`) |

> **Security Note**: The `Password` field uses `json:"-"` to prevent password hashes from being exposed in API responses.

---

### UserSettings
User preferences and configuration.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `UserID` | string | Foreign key to User (unique) |
| `DefaultCurrency` | string | Default currency code (default: "USD") |
| `DateFormat` | string | Date format string (default: "2006-01-02") |
| `Theme` | string | UI theme (default: "light") |
| `Language` | string | Language code (default: "en") |
| `Notifications` | bool | Enable notifications (default: true) |

---

### Session
Refresh tokens for JWT authentication.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `UserID` | string | Foreign key to User |
| `TokenHash` | string | Hashed refresh token (never exposed) |
| `ExpiresAt` | time.Time | Token expiration time |
| `RevokedAt` | *time.Time | When session was revoked |

---

## Financial Models

### Account
Represents a financial account (bank, cash, credit card, etc.).

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `Name` | string | Account name |
| `Type` | string | Account type (e.g., "Bank", "Cash", "Credit Card") |
| `Balance` | float64 | Current balance |
| `Currency` | string | Currency code (default: "USD") |

---

### Category
Income/expense categories.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `Name` | string | Category name |
| `Icon` | string | Icon identifier |
| `Color` | string | Color code (hex) |
| `Type` | string | "Income" or "Expense" |

---

### Tag
Flexible categorization of transactions (many-to-many).

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `Name` | string | Tag name |
| `Color` | string | Hex color code |

---

### TransactionTag
Join table for Transaction ↔ Tag relationship.

| Field | Type | Description |
|-------|------|-------------|
| `TransactionID` | string | Composite primary key |
| `TagID` | string | Composite primary key |

---

### Transaction
Core model for tracking income and expenses.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `AccountID` | string | Foreign key to Account |
| `CategoryID` | string | Foreign key to Category |
| `TransferID` | *string | Optional link to Transfer |
| `RecurringPaymentID` | *string | Optional link to RecurringPayment |
| `DebtID` | *string | Optional link to Debt |
| `DebtInstallmentID` | *string | Optional link to DebtInstallment |
| `Amount` | float64 | Transaction amount |
| `Description` | string | Description |
| `Date` | time.Time | Transaction date |
| `Type` | string | "Income" or "Expense" |

---

### Transfer
Money movement between accounts.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `FromAccountID` | string | Source account |
| `ToAccountID` | string | Destination account |
| `Amount` | float64 | Transfer amount |
| `Description` | string | Description |
| `Date` | time.Time | Transfer date |

---

### RecurringPayment
Scheduled recurring transactions.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `AccountID` | *string | Optional linked account |
| `CategoryID` | string | Foreign key to Category |
| `Amount` | float64 | Payment amount |
| `Frequency` | string | "Monthly", "Weekly", etc. |
| `StartDate` | time.Time | When recurrence starts |
| `NextDate` | time.Time | Next occurrence date |
| `IsActive` | bool | Whether the recurrence is active |

---

### RecurringInstance
Tracks generated transactions from recurring payments.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `RecurringPaymentID` | string | Foreign key to RecurringPayment |
| `TransactionID` | *string | Generated transaction |
| `DueDate` | time.Time | When this instance is due |
| `Status` | string | "pending", "processed", or "skipped" |

---

### Debt
Money lent or borrowed.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `Name` | string | Debt name/description |
| `TotalAmount` | float64 | Original amount |
| `Remaining` | float64 | Amount still owed |
| `DueDate` | time.Time | When debt is due |
| `Type` | string | "Lent" or "Borrowed" |

---

### DebtInstallment
Individual payments toward a debt.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `DebtID` | string | Foreign key to Debt |
| `Amount` | float64 | Payment amount |
| `PaidDate` | time.Time | When payment was made |

---

## Planning Models

### Wishlist
Savings goals for desired purchases.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `Name` | string | Item name |
| `Description` | string | Description |
| `TargetAmount` | float64 | Cost of the item |
| `CurrentSaved` | float64 | Amount saved so far |
| `Need` | int | Priority 1-6 |
| `Want` | int | Desire level 1-6 |
| `Productivity` | int | Utility score 1-6 |
| `URL` | string | Product link |
| `Status` | string | "Active", "Purchased", "Archived" |
| `DesiredDate` | *time.Time | Target purchase date |

#### Virtual Fields (Calculated at runtime)
- `FinancialImpact`: Calculated financial impact score
- `Score`: Calculated priority score

---

### Budget
Spending limits per category or globally.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `CategoryID` | *string | Optional category (null = global budget) |
| `Amount` | float64 | Budget limit |
| `Period` | string | "Monthly", "Weekly", "Yearly", "One-Time" |
| `StartDate` | time.Time | Budget start date |
| `EndDate` | *time.Time | Optional budget end date |
| `IsStrict` | bool | Whether exceeding is blocked |

---

### SavingsGoal
Long-term savings targets.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `Name` | string | Goal name |
| `TargetAmount` | float64 | Savings target |
| `CurrentAmount` | float64 | Current progress |
| `AccountID` | *string | Optional dedicated account |
| `TargetDate` | *time.Time | Optional deadline |

---

## Audit & Storage Models

### Attachment
Receipt images and files linked to transactions.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `TransactionID` | *string | Optional linked transaction |
| `FileName` | string | Original file name |
| `FilePath` | string | Relative path in storage |
| `MimeType` | string | MIME type (e.g., "image/jpeg") |
| `Size` | int64 | File size in bytes |

---

### AuditLog
Change history for debugging and compliance.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Primary key |
| `TenantID` | string | Foreign key to Tenant |
| `UserID` | string | Foreign key to User |
| `Action` | string | "create", "update", "delete" |
| `EntityType` | string | "transaction", "account", etc. |
| `EntityID` | string | ID of the affected entity |
| `OldValue` | string | JSON of previous state |
| `NewValue` | string | JSON of new state |
| `IPAddress` | string | Client IP address |
| `UserAgent` | string | Client user agent |

---

## Entity Relationships

```
Tenant (1) ──── (N) User
Tenant (1) ──── (N) Account
Tenant (1) ──── (N) Category
Tenant (1) ──── (N) Tag
Tenant (1) ──── (N) Transaction
Tenant (1) ──── (N) Transfer
Tenant (1) ──── (N) RecurringPayment
Tenant (1) ──── (N) Debt
Tenant (1) ──── (N) Wishlist
Tenant (1) ──── (N) Budget
Tenant (1) ──── (N) SavingsGoal
Tenant (1) ──── (N) AuditLog

User (1) ──── (1) UserSettings
User (1) ──── (N) Session

Account (1) ──── (N) Transaction
Category (1) ──── (N) Transaction
Tag (N) <───── (N) Transaction (via TransactionTag)
Transfer (1) ──── (N) Transaction
RecurringPayment (1) ──── (N) RecurringInstance
Debt (1) ──── (N) DebtInstallment
Transaction (1) ──── (N) Attachment
```

---

## Usage

Models are migrated automatically on startup in `main.go`:

```go
db.AutoMigrate(
    &models.Tenant{},
    &models.User{},
    &models.UserSettings{},
    &models.Session{},
    // ... other models
)
```

---

## Best Practices

1. **Always include TenantID** for multi-tenant isolation
2. **Use `json:"-"`** for sensitive fields like passwords
3. **Use pointers** for nullable fields and optional relationships
4. **Leverage GORM hooks** (BeforeCreate, AfterUpdate) for computed fields
5. **Use `gorm:"-"`** for virtual fields that shouldn't be persisted