# Controllers Layer

Controllers are the HTTP handlers that process incoming requests and return responses. They act as the entry point for the API.

## Overview

The controllers layer follows a constructor pattern for dependency injection, making it easy to test and maintain.

---

## BaseController

`base_controller.go` provides standardized HTTP response utilities that all controllers embed.

### Methods

| Method | Description |
|--------|-------------|
| `JSON(w, status, data)` | Responds with a JSON payload |
| `Error(w, status, message)` | Responds with a standardized error object |
| `Bind(r, v)` | Parses JSON request body into a struct |
| `GetPaginationParams(r)` | Extracts pagination parameters from query string |
| `PaginatedResponse(w, p)` | Responds with paginated results |

### Example Usage

```go
// Success response
c.JSON(w, http.StatusOK, user)

// Error response
c.Error(w, http.StatusBadRequest, "invalid request body")

// Bind request body
var user models.User
if err := c.Bind(r, &user); err != nil {
    c.Error(w, http.StatusBadRequest, "invalid request body")
    return
}

// Pagination
pagination := c.GetPaginationParams(r)
// Returns &Pagination{Limit: 10, Page: 1, Sort: "id DESC"}
```

---

## AuthController

`auth_controller.go` handles authentication-related endpoints.

### Dependencies

- `BaseController` - Embedded for response utilities
- `AuthService` - Business logic for authentication

### Constructor

```go
func NewAuthController(db *gorm.DB) *AuthController
```

Creates a new controller with all dependencies initialized.

### Endpoints

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/api/auth/register` | `Register` | Register a new user |
| POST | `/api/auth/login` | `Login` | Authenticate a user |

### Register Handler

```go
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request)
```

- Accepts JSON body with user details
- Validates and creates a new user
- Password is hashed by the service layer
- Returns `201 Created` on success

**Request Body:**
```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123",
    "tenant_id": "tenant_ulid"
}
```

**Response:**
```json
{
    "id": "01HXYZ...",
    "name": "John Doe",
    "email": "john@example.com",
    "tenant_id": "tenant_ulid",
    "created_at": "2024-01-01T00:00:00Z"
}
```

### Login Handler

```go
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request)
```

- Accepts email and password
- Validates credentials
- Returns user data on success

**Request Body:**
```json
{
    "email": "john@example.com",
    "password": "secret123"
}
```

**Response:**
```json
{
    "id": "01HXYZ...",
    "name": "John Doe",
    "email": "john@example.com"
}
```

---

## Creating a New Controller

### Template

```go
package controllers

import (
    "gosvelte/app/models"
    "gosvelte/app/services"
    "net/http"
    "gorm.io/gorm"
)

type MyController struct {
    BaseController
    Service *services.MyService
}

func NewMyController(db *gorm.DB) *MyController {
    return &MyController{
        Service: services.NewMyService(db),
    }
}

func (c *MyController) MyHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var input models.MyModel
    if err := c.Bind(r, &input); err != nil {
        c.Error(w, http.StatusBadRequest, "invalid request body")
        return
    }

    // 2. Call service
    result, err := c.Service.DoSomething(&input)
    if err != nil {
        c.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    // 3. Return response
    c.JSON(w, http.StatusOK, result)
}
```

### Steps

1. Create a new file in `app/controllers/`
2. Embed `BaseController` for response utilities
3. Define your controller struct with service dependencies
4. Create a constructor function that initializes dependencies
5. Implement handler methods
6. Register routes in `app/routes.go`

---

## HTTP Status Codes

Use proper status codes for different outcomes:

| Status | Code | Usage |
|--------|------|-------|
| OK | 200 | Successful GET/PUT |
| Created | 201 | Successful POST |
| No Content | 204 | Successful DELETE |
| Bad Request | 400 | Invalid input |
| Unauthorized | 401 | Authentication required |
| Forbidden | 403 | No permission |
| Not Found | 404 | Resource doesn't exist |
| Conflict | 409 | Duplicate resource |
| Internal Server Error | 500 | Server-side error |

---

## Best Practices

1. **Keep controllers thin** - They should only parse requests and format responses
2. **Delegate to services** - Business logic belongs in the service layer
3. **Validate input** - Use `Bind()` and check for required fields
4. **Return appropriate status codes** - Use the table above
5. **Don't expose internal errors** - Return generic messages for 500 errors
6. **Use dependency injection** - Makes testing easier