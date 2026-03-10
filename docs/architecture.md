# Architecture Overview

The GoSvelte application follows a strict **MVC + Service + Repository** architecture pattern.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        HTTP Request                              │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Routes (routes.go)                          │
│                    URL → Controller mapping                      │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Controllers Layer                             │
│              Parse request → Call service → Respond              │
│                                                                  │
│  ┌──────────────────┐  ┌──────────────────┐                     │
│  │  BaseController  │  │  AuthController  │  ...                 │
│  └──────────────────┘  └──────────────────┘                     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Services Layer                               │
│              Business logic + Transaction management             │
│                                                                  │
│  ┌──────────────────┐  ┌──────────────────┐                     │
│  │   BaseService    │  │   AuthService   │                     │
│  │                  │  │  FinanceService │  ...                 │
│  └──────────────────┘  └──────────────────┘                     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Repositories Layer                             │
│                    Database operations                           │
│                                                                  │
│  ┌──────────────────────────────────────────┐                   │
│  │  BaseRepository[T] (Generic CRUD)         │                   │
│  │  - FindByID, FindAll, PaginatedFind       │                   │
│  │  - Create, Update, Delete                │                   │
│  │  - FilterScope, SearchScope              │                   │
│  └──────────────────────────────────────────┘                   │
│  ┌──────────────────┐  ┌──────────────────┐                     │
│  │ UserRepository   │  │ AccountRepo     │  ...                 │
│  └──────────────────┘  └──────────────────┘                     │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Models Layer                                │
│                    GORM structs + relations                      │
│                                                                  │
│  Tenant, User, Session, UserSettings, Account, Category, Tag,    │
│  Transaction, Transfer, RecurringPayment, Debt, Wishlist,       │
│  Budget, SavingsGoal, Attachment, AuditLog, ...                 │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      PostgreSQL Database                         │
└─────────────────────────────────────────────────────────────────┘
```

---

## Directory Structure

```
app/
├── controllers/          # HTTP handlers
│   ├── base_controller.go    # Response utilities
│   ├── auth_controller.go    # Authentication endpoints
│   └── README.md              # (moved to docs/controllers.md)
│
├── services/             # Business logic
│   ├── base_service.go       # Shared service utilities
│   ├── auth_service.go       # Authentication logic
│   ├── finance_service.go    # Financial operations
│   └── README.md              # (moved to docs/services.md)
│
├── repositories/         # Data access
│   ├── base_repository.go    # Generic CRUD operations
│   ├── user_repository.go    # User-specific queries
│   ├── *_repository.go       # Entity repositories
│   └── README.md              # (moved to docs/repositories.md)
│
├── models/               # Database models
│   ├── finance.go            # All GORM models
│   └── README.md              # (moved to docs/models.md)
│
└── routes.go             # API route definitions

docs/
├── models.md             # Models documentation
├── controllers.md        # Controllers documentation
├── services.md           # Services documentation
├── repositories.md       # Repositories documentation
├── architecture.md       # This file
└── services_plan.txt     # Service architecture plan
```

---

## Request Flow

### 1. Route Matching (`routes.go`)

```go
mux.HandleFunc("POST /api/auth/register", authController.Register)
```

### 2. Controller Processing (`controllers/auth_controller.go`)

```go
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := c.Bind(r, &user); err != nil {
        c.Error(w, http.StatusBadRequest, "invalid request body")
        return
    }
    if err := c.Service.Register(&user); err != nil {
        c.Error(w, http.StatusInternalServerError, err.Error())
        return
    }
    c.JSON(w, http.StatusCreated, user)
}
```

### 3. Service Business Logic (`services/auth_service.go`)

```go
func (s *AuthService) Register(user *models.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    return s.Repo.Create(user)
}
```

### 4. Repository Data Access (`repositories/user_repository.go`)

```go
func (r *UserRepository) Create(user *models.User) error {
    return r.DB.Create(user).Error
}
```

---

## Dependency Injection

The application uses constructor-based dependency injection:

```go
// In routes.go
func RegisterRoutes(db *gorm.DB) *http.ServeMux {
    // Controller creates its own dependencies
    authController := controllers.NewAuthController(db)
    // ...
}

// In auth_controller.go
func NewAuthController(db *gorm.DB) *AuthController {
    repo := &repositories.UserRepository{
        BaseRepository: repositories.BaseRepository[models.User]{DB: db},
    }
    service := &services.AuthService{
        BaseService: services.BaseService{DB: db},
        Repo:        repo,
    }
    return &AuthController{Service: service}
}
```

---

## API Endpoints

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/api/auth/register` | AuthController.Register | User registration |
| POST | `/api/auth/login` | AuthController.Login | User authentication |
| GET/POST | `/static/*` | FileServer | Static assets |
| GET | `/*` | SPA handler | Serves frontend |

---

## Multi-Tenancy

All models (except `Tenant`) include a `TenantID` field for data isolation:

```go
type User struct {
    ID       string `json:"id"`
    TenantID string `json:"tenant_id"`  // Multi-tenant isolation
    // ...
}
```

Services should always filter by tenant:

```go
func (s *MyService) GetItems(tenantID string) ([]models.Item, error) {
    return s.Repo.FindAll(s.Repo.FilterScope("tenant_id", "=", tenantID))
}
```

---

## Adding a New Feature

1. **Define the model** in `app/models/finance.go`
2. **Create repository** in `app/repositories/` (extend BaseRepository)
3. **Create service** in `app/services/` (embed BaseService)
4. **Create controller** in `app/controllers/` (embed BaseController)
5. **Register routes** in `app/routes.go`
6. **Add migration** in `main.go`

---

## Testing Strategy

Each layer can be tested independently:

| Layer | Test Type | Approach |
|-------|-----------|----------|
| Models | Unit | Test struct methods and computed fields |
| Repositories | Integration | Test with a test database |
| Services | Unit | Mock repositories |
| Controllers | HTTP | Mock services |

---

## Security Considerations

1. **Password fields** use `json:"-"` to prevent exposure
2. **Tenant isolation** is enforced at the service layer
3. **Soft deletes** prevent accidental data loss
4. **Input validation** happens at the controller layer
5. **SQL injection** is prevented by GORM parameterized queries
6. **Session tokens** are hashed before storage

---

## Related Documentation

- [Models](models.md) - Database schema and relationships
- [Controllers](controllers.md) - HTTP handlers and routing
- [Services](services.md) - Business logic layer
- [Repositories](repositories.md) - Data access layer
- [Services Plan](services_plan.txt) - Planned service implementations