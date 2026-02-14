# --- Stage 1: Build Frontend ---
FROM node:20-alpine AS frontend-builder
WORKDIR /build/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# --- Stage 2: Build Backend ---
FROM golang:1.24-alpine AS backend-builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the built frontend from Stage 1
COPY --from=frontend-builder /build/static ./static
# Build the Go binary (CGO is NOT needed for Postgres)
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# --- Stage 3: Runtime ---
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
# Copy the binary and static files
COPY --from=backend-builder /build/server .
COPY --from=backend-builder /build/static ./static

EXPOSE 8080
CMD ["./server"]
