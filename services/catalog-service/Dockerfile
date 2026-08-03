# ==========================================
# 1. Build Stage
# ==========================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Copy go.work and all module manifests in the workspace
COPY go.work* ./
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY services/user-service/go.mod services/user-service/go.sum* ./services/user-service/
COPY services/catalog-service/go.mod services/catalog-service/go.sum* ./services/catalog-service/

# Download dependencies for catalog-service
WORKDIR /app/services/catalog-service
RUN go mod download

# Copy source code for pkg and catalog-service
WORKDIR /app
COPY pkg/ ./pkg/
COPY services/catalog-service/ ./services/catalog-service/

# Build binary
WORKDIR /app/services/catalog-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/catalog-service ./cmd/main.go

# ==========================================
# 2. Final Runtime Stage
# ==========================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/catalog-service .

EXPOSE 8081
ENTRYPOINT ["./catalog-service"]