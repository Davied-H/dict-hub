# =============================================================================
# Dict-Hub Dockerfile - Multi-stage Build with Embedded Static Files
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Frontend Build (Node.js)
# -----------------------------------------------------------------------------
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Install dependencies first (better layer caching)
COPY frontend/package*.json ./
RUN npm ci --silent

# Copy source and build
COPY frontend/ ./
RUN npm run build

# -----------------------------------------------------------------------------
# Stage 2: Backend Build (Go) - Embed frontend static files
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install build dependencies (CGO required for SQLite)
RUN apk add --no-cache gcc musl-dev

# Download dependencies first (better layer caching)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Copy frontend build output to web/dist for go:embed
# Note: vite outputs to ../backend/web/dist relative to frontend dir
COPY --from=frontend-builder /app/backend/web/dist ./web/dist

# Build with embedded static files
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# -----------------------------------------------------------------------------
# Stage 3: Production Image (Alpine)
# -----------------------------------------------------------------------------
FROM alpine:3.20

# Install runtime dependencies
# - libc6-compat: Required for CGO/SQLite compatibility
# - ca-certificates: For HTTPS requests
# - tzdata: For timezone support
RUN apk add --no-cache ca-certificates tzdata libc6-compat

# Set timezone
ENV TZ=Asia/Shanghai

WORKDIR /app

# Copy binary from backend builder (includes embedded frontend)
COPY --from=backend-builder /app/server .

# Copy default config
COPY backend/configs/config.yaml ./configs/

# Create data directories with proper permissions
# Using 777 to ensure compatibility with various unraid user configurations
RUN mkdir -p /app/data /app/dicts/source /app/dicts/sound && \
    chmod -R 777 /app/data /app/dicts

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./server"]
