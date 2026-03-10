# Repositories Layer

The repositories layer handles all database operations using Go generics for code reuse.

## Overview

The repository pattern abstracts database access, allowing services to focus on business logic without knowing SQL details.

---

## BaseRepository

`base_repository.go` provides a generic CRUD interface that works with any model type using Go generics.

### Type Definition

```go
type BaseRepository[T any] struct {
    DB *gorm.DB
}
```

### CRUD Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `FindByID` | `(id string) (*T, error)` | Find a single record by ID |
| `FindAll` | `(scopes ...func(*gorm.DB) *gorm.DB) ([]T, error)` | Find all records with optional filters |
| `Create` | `(item *T) error` | Insert a new record |
| `Update` | `(item *T) error` | Save changes to a record |
| `Delete` | `(item *T) error` | Soft delete a record |

### Pagination Method

| Method | Signature | Description |
|--------|-----------|-------------|
| `PaginatedFind` | `(pagination *Pagination, scopes ...func(*gorm.DB) *gorm.DB) (*Pagination, error)` | Find records with pagination |

### Scope Methods

Scopes are reusable query modifiers that can be chained:

| Method | Description |
|--------|-------------|
| `FilterScope(column, operator, value)` | Creates a WHERE clause filter |
| `SearchScope(query, columns...)` | Creates a LIKE search across multiple columns |

---

## Usage Examples

### FilterScope

```go
repo := &BaseRepository[models.User]{DB: db}

// Filter by tenant
tenantFilter := repo.FilterScope("tenant_id", "=", "tenant_123")

users, err := repo.FindAll(tenantFilter)
```

### SearchScope

```go
repo := &BaseRepository[models.User]{DB: db}

// Search in name and email
searchScope := repo.SearchScope("john", "name", "email")

users, err := repo.FindAll(searchScope)
// Generates: WHERE name LIKE '%john%' OR email LIKE '%john%'
```

### Chaining Scopes

```go
// Combine multiple scopes
results, err := repo.FindAll(
    repo.FilterScope("tenant_id", "=", tenantID),
    repo.FilterScope("is_active", "=", true),
    repo.SearchScope(search, "name", "description"),
)
```

---

## Pagination

The `Pagination` struct handles result paging.

### Struct Definition

```go
type Pagination struct {
    Limit      int         `json:"limit"`       // Items per page
    Page       int         `json:"page"`        // Current page
    Sort       string      `json:"sort"`        // Sort order (e.g., "id desc")
    TotalRows  int64       `json:"total_rows"`  // Total matching records
    TotalPages int         `json:"total_pages"` // Calculated total pages
    Rows       interface{} `json:"rows"`        // Result data
}
```

### Usage

```go
pagination := &Pagination{
    Limit: 20,
    Page: 1,
    Sort: "created_at DESC",
}

result, err := repo.PaginatedFind(pagination, tenantFilter)
// result.Rows contains the data
// result.TotalPages is automatically calculated
```

### Helper Methods

| Method | Description |
|--------|-------------|
| `GetOffset()` | Calculates offset for LIMIT clause |
| `GetLimit()` | Returns limit (default: 10) |
| `GetPage()` | Returns page number (default: 1) |
| `GetSort()` | Returns sort clause (default: "id DESC") |

---

## UserRepository

`user_repository.go` extends `BaseRepository` with custom queries.

### Struct

```go
type UserRepository struct {
    BaseRepository[models.User]
}
```

### Custom Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `FindByEmail` | `(email string) (*models.User, error)` | Find user by email address |

### Usage

```go
// Initialize
userRepo := &repositories.UserRepository{
    BaseRepository: repositories.BaseRepository[models.User]{DB: db},
}

// Use base methods
users, err := userRepo.FindAll()

// Use custom methods
user, err := userRepo.FindByEmail("john@example.com")
```

---

## ULID Generation

The repository includes a ULID generation utility:

```go
func GenerateULID() string
```

ULIDs are:
- Universally unique (like UUIDs)
- Sortable by timestamp
- URL-safe (base32 encoded)

---

## Repository Files

| File | Description |
|------|-------------|
| `base_repository.go` | Generic CRUD with pagination and scopes |
| `user_repository.go` | User-specific queries |
| `account_repository.go` | Account repository stub |
| `category_repository.go` | Category repository stub |
| `transaction_repository.go` | Transaction repository stub |
| `transfer_repository.go` | Transfer repository stub |
| `recurring_payment_repository.go` | Recurring payment repository stub |
| `debt_repository.go` | Debt repository stub |
| `debt_installment_repository.go` | Debt installment repository stub |
| `wishlist_repository.go` | Wishlist repository stub |
| `budget_repository.go` | Budget repository stub |
| `savings_goal_repository.go` | Savings goal repository stub |
| `tenant_repository.go` | Tenant repository stub |

> **Note**: Stub files are minimal structs that inherit all base functionality. Add custom methods as needed.

---

## Creating a Custom Repository

For entities that need custom queries beyond CRUD:

### Template

```go
package repositories

import "gosvelte/app/models"

type AccountRepository struct {
    BaseRepository[models.Account]
}

// Custom method
func (r *AccountRepository) FindByTenantID(tenantID string) ([]models.Account, error) {
    var accounts []models.Account
    err := r.DB.Where("tenant_id = ?", tenantID).Find(&accounts).Error
    return accounts, err
}

// Method with joins
func (r *AccountRepository) FindWithTransactions(accountID string) (*models.Account, error) {
    var account models.Account
    err := r.DB.Preload("Transactions").First(&account, "id = ?", accountID).Error
    return &account, err
}
```

---

## Best Practices

1. **Use BaseRepository for simple CRUD** - Don't duplicate basic operations
2. **Create custom methods for complex queries** - Joins, aggregations, etc.
3. **Return pointers** - Methods return `*T` to allow nil checks
4. **Use scopes for reusable filters** - Don't repeat WHERE clauses
5. **Let GORM handle errors** - Return errors directly to the service layer
6. **Always filter by tenant** - Add `FilterScope("tenant_id", "=", tenantID)` to queries