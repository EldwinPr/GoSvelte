# Services Layer

The services layer contains business logic and transaction management. Controllers delegate complex operations to services.

## Overview

Services encapsulate:
- Business rules and validations
- Database transactions
- Cross-entity operations
- External integrations

---

## BaseService

`base_service.go` provides shared functionality for all services.

### Embedded Fields

| Field | Type | Description |
|-------|------|-------------|
| `DB` | `*gorm.DB` | Database connection for queries |

### Methods

| Method | Description |
|--------|-------------|
| `Scope(fn)` | Injects DB scopes for reusable query modifiers |
| `Transaction(fn)` | Wraps operations in a database transaction |

### Scope Example

```go
// Define a reusable scope
activeScope := func(db *gorm.DB) *gorm.DB {
    return db.Where("is_active = ?", true)
}

// Use in service
result := s.Scope(activeScope).Find(&items)
```

### Transaction Example

```go
err := s.Transaction(func(tx *gorm.DB) error {
    // All operations within this function use the same transaction
    if err := tx.Create(&user).Error; err != nil {
        return err // Will rollback
    }
    if err := tx.Create(&profile).Error; err != nil {
        return err // Will rollback
    }
    return nil // Will commit
})
```

---

## Service Implementation Standards

To ensure the GoSvelte application remains robust, atomic, and secure, all service implementations must adhere to these standards.

### 1. Mandatory Tenancy Isolation
Every service method (except for initial Login/Registration) **must** accept `tenantID string` as its first argument.
- **Rule:** Services must verify ownership of every resource being accessed.
- **Example:** `WHERE id = accountID AND tenant_id = tenantID`.

### 2. Transactional Atomicity (The "All or Nothing" Rule)
Any operation that modifies more than one table or performs multiple steps **must** be wrapped in a database transaction.
- **Pattern:** Use the `s.Transaction(func(tx *gorm.DB) error { ... })` helper from `BaseService`.
- **Failure:** If any step fails inside the function, all changes must be automatically rolled back.

### 3. Service-to-Service Orchestration
Services should delegate complex domain logic to other services to avoid logic duplication.
- **Dependency Rule:** Services may depend on other services for complex operations (e.g., `WishlistService` calling `TransactionService`).
- **Transaction Propagation:** When Service A calls Service B, it must pass its current `*gorm.DB` (the transaction handle) to ensure both services operate within the same atomic unit.

### 4. Domain-Driven Validation
Services are the gatekeepers of business rules and are responsible for enforcing domain logic before any database changes occur.
- **Pre-conditions:** Check for sufficient balances, valid dates, and logical consistency.
- **Error Mapping:** Services should return specific domain errors (e.g., `ErrInsufficientFunds`) that controllers can then map to appropriate HTTP status codes.

### 5. Separation of Concerns
- **Repositories:** Responsible for *how* to access data (SQL, Joins, Filters, Scopes).
- **Services:** Responsible for *what* to do with that data (Business rules, Calculations, Orchestration).

### 6. Concurrency Safety (Locking)
Critical financial operations (like updating account balances) must use row-level locking to prevent race conditions.
- **Pattern:** Use `FOR UPDATE` locks when fetching accounts that are about to have their balances modified.

---

## AuthService
`auth_service.go` handles user authentication logic.


### Dependencies

- `BaseService` - For database access
- `UserRepository` - For user data operations

### Methods

| Method | Description |
|--------|-------------|
| `Register(user)` | Creates a new user with hashed password |
| `Login(email, password)` | Authenticates user credentials |

### Register Flow

```go
func (s *AuthService) Register(user *models.User) error {
    // 1. Hash the password with bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // 2. Replace plain password with hash
    user.Password = string(hashedPassword)

    // 3. Create user in database
    return s.Repo.Create(user)
}
```

### Login Flow

```go
func (s *AuthService) Login(email, password string) (*models.User, error) {
    // 1. Find user by email
    user, err := s.Repo.FindByEmail(email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // 2. Compare password hash
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // 3. Return user (password already excluded from JSON)
    return user, nil
}
```

> **Security Note**: Error messages are intentionally vague ("invalid credentials") to prevent email enumeration attacks.

---

## FinanceService

`finance_service.go` handles financial transactions with atomic operations.

### Dependencies

- `BaseService` - For database access and transactions
- `AccountRepo` - For account operations
- `TransactionRepo` - For transaction operations

### Constructor

```go
func NewFinanceService(db *gorm.DB) *FinanceService
```

### Methods

| Method | Description |
|--------|-------------|
| `CreateTransaction(txData)` | Creates a transaction and updates account balance atomically |
| `ReconcileAccount(tenantID, accountID)` | Recalculates account balance from transaction history |

### CreateTransaction Flow

This method demonstrates a critical transaction pattern:

```go
func (s *FinanceService) CreateTransaction(txData *models.Transaction) error {
    return s.Transaction(func(tx *gorm.DB) error {
        // 1. Lock the account for concurrent safety
        var account models.Account
        if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&account, txData.AccountID).Error; err != nil {
            return err
        }

        // 2. Validate tenant ownership
        if account.TenantID != txData.TenantID {
            return errors.New("unauthorized account access")
        }

        // 3. Update balance based on transaction type
        if txData.Type == "Income" {
            account.Balance += txData.Amount
        } else if txData.Type == "Expense" {
            account.Balance -= txData.Amount
        }

        // 4. Save account
        if err := tx.Save(&account).Error; err != nil {
            return err
        }

        // 5. Create transaction record
        return tx.Create(txData).Error
    })
}
```

**Key Points:**
- Uses `FOR UPDATE` lock to prevent race conditions
- Validates tenant ownership for security
- Updates balance and creates transaction atomically
- All changes rollback if any step fails

### ReconcileAccount Flow

Recalculates account balance from transaction history:

```go
func (s *FinanceService) ReconcileAccount(tenantID, accountID uint) (float64, error) {
    // 1. Sum all income transactions
    var income float64
    s.DB.Model(&models.Transaction{}).
        Where("tenant_id = ? AND account_id = ? AND type = ?", tenantID, accountID, "Income").
        Select("COALESCE(SUM(amount), 0)").
        Scan(&income)

    // 2. Sum all expense transactions
    var expense float64
    s.DB.Model(&models.Transaction{}).
        Where("tenant_id = ? AND account_id = ? AND type = ?", tenantID, accountID, "Expense").
        Select("COALESCE(SUM(amount), 0)").
        Scan(&expense)

    // 3. Calculate balance
    total := income - expense

    // 4. Update account
    err := s.DB.Model(&models.Account{}).
        Where("id = ? AND tenant_id = ?", accountID, tenantID).
        Update("balance", total).Error

    return total, err
}
```

---

## Service Architecture Plan

Based on `docs/services_plan.txt`, the following services are planned:

| Service | Purpose |
|---------|---------|
| `AuthService` | User registration, login, session management |
| `AccountService` | Account CRUD, dashboard summary, reconciliation |
| `TransactionService` | Income/expense recording, ledger queries |
| `TransferService` | Account-to-account transfers |
| `DebtService` | Debt tracking, installment payments |
| `RecurringPaymentService` | Scheduled transactions |
| `BudgetService` | Budget limits, performance tracking |
| `WishlistService` | Purchase planning, smart scoring |
| `SavingsGoalService` | Long-term savings targets |
| `ReportingService` | Analytics, spending breakdown, cash flow |

---

## Creating a New Service

### Template

```go
package services

import (
    "gorm.io/gorm"
    "gosvelte/app/models"
    "gosvelte/app/repositories"
)

type MyService struct {
    BaseService
    MyRepo *repositories.BaseRepository[models.MyModel]
}

func NewMyService(db *gorm.DB) *MyService {
    return &MyService{
        BaseService: BaseService{DB: db},
        MyRepo:      &repositories.BaseRepository[models.MyModel]{DB: db},
    }
}

func (s *MyService) DoSomething(input *models.MyModel) error {
    // Business logic here
    return s.MyRepo.Create(input)
}
```

---

## Best Practices

1. **Use transactions for multi-step operations** - Wrap related operations with `s.Transaction()`
2. **Validate business rules** - Services should enforce domain constraints
3. **Return meaningful errors** - Use specific error messages for debugging
4. **Keep services focused** - One service per domain concept
5. **Inject repositories** - Don't create repositories inside methods
6. **Lock rows when needed** - Use `FOR UPDATE` for concurrent operations
7. **Always filter by tenant** - Enforce data isolation