# GoSvelte ERP Framework

A professional-grade, microservice-ready boilerplate using Go (MVC+Service) and Svelte 5 (Vite + Tailwind CSS v4).

## 🚀 Architecture

### Backend (Go)
- **Pattern**: Strict MVC + Service + Repository.
- **Models**: GORM structs in `app/models`.
- **Repositories**: Generic Base Repository in `app/repositories` to handle CRUD without boilerplate.
- **Services**: Business logic and transaction management in `app/services`.
- **Controllers**: HTTP handlers in `app/controllers`. Uses a Constructor pattern for Dependency Injection.
- **Routing**: Centralized in `app/routes.go` using a `/api` prefix to isolate backend logic.

### Frontend (Svelte)
- **Framework**: Svelte 5 + TypeScript.
- **Styling**: Tailwind CSS v4 (using the native Vite plugin).
- **Routing**: Hash-based routing via `svelte-spa-router`.
- **API Client**: Standardized fetch wrapper in `web/src/lib/api.ts`.

---

## 🛠️ Development Setup

### Prerequisites
- Go 1.24+
- Node.js 20+
- Docker Desktop (for PostgreSQL)

### Local Development (Recommended)
1. **Start the Database**:
   ```bash
   docker compose up db -d
   ```
2. **Start the Backend (with hot-reload)**:
   ```bash
   air
   ```
3. **Start the Frontend**:
   ```bash
   cd web
   npm install
   npm run dev
   ```

### Production Build
To build the entire application into a single containerized image:
```bash
docker compose up --build
```

---

## 📁 Project Structure
```text
├── app/
│   ├── controllers/    # HTTP Handlers
│   ├── models/         # Database GORM Structs
│   ├── repositories/   # Data Access Layer (Generic)
│   ├── services/       # Business Logic Layer
│   └── routes.go       # API Route Definitions
├── web/                # Svelte 5 Frontend
│   ├── src/pages/      # UI Views
│   └── src/lib/        # Shared Frontend Logic
├── static/             # Compiled Frontend Assets
├── main.go             # Entry Point
└── docker-compose.yml  # Infrastructure Orchestration
```

## 🔒 Security Notes
- Passwords are hidden from JSON output by default (`json:"-"` tag).
- The Go `BaseController` provides standardized error and success response formats.
