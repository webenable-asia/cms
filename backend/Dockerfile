# backend/Dockerfile
FROM golang:1.24-alpine

# Install git (required for go mod download)
RUN apk add --no-cache git

WORKDIR /app

# Install air for live reloading
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

# Use air for live-reloading
CMD ["air"]
