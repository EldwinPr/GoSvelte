# Makefile for GoSvelte ERP Framework

.PHONY: dev-backend dev-frontend dev build-frontend build-backend build clean install

# --- Development ---

# Run Go backend with Air for hot-reloading
dev-backend:
	air

# Run Svelte frontend with Vite
dev-frontend:
	cd web && npm run dev

# Run both simultaneously (requires 'npx concurrently' or similar, or just run in two terminals)
dev:
	npx concurrently "make dev-backend" "make dev-frontend"

# --- Production Build ---

# Build Svelte frontend and output to /static
build-frontend:
	cd web && npm run build

# Build Go backend
build-backend:
	go build -o bin/server.exe main.go

# Full build: Frontend first, then Backend
build: build-frontend build-backend

# --- Utilities ---

install:
	go mod download
	cd web && npm install

clean:
	rm -rf static/*
	rm -rf bin/*
	rm -rf web/dist
