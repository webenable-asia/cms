# Multi-stage build for optimized Go application
ARG GO_ENV=development

# ===== Dependencies Stage =====
FROM docker.io/golang:1.24-alpine AS deps
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# ===== Builder Stage (Production only) =====
FROM docker.io/golang:1.24-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy dependencies
COPY --from=deps /go/pkg /go/pkg
COPY go.mod go.sum ./
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

# ===== Production Stage =====
FROM docker.io/alpine:3.22 AS production
WORKDIR /app

# Install CA certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -D appuser

# Create uploads directory with proper permissions
RUN mkdir -p /app/uploads && \
    chown -R appuser:appgroup /app

# Copy the binary
COPY --from=builder /app/main .

EXPOSE 8080
USER appuser

CMD ["./main"]

# ===== Development Stage =====
FROM docker.io/golang:1.24-alpine AS development
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Install air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy dependencies from deps stage
COPY --from=deps /go/pkg /go/pkg
COPY go.mod go.sum ./

# Copy source code
COPY . .

EXPOSE 8080

CMD ["air"]

# ===== Final Stage =====
FROM ${GO_ENV} AS final
