# ==========================================
# 1. Build Stage
# ==========================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Copy the workspace and every module manifest so the full module graph
# can be resolved. Source for unneeded services is NOT copied below, so
# their code is never compiled.
COPY go.work go.work.sum* ./
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY services/user-service/go.mod services/user-service/go.sum* ./services/user-service/
COPY services/catalog-service/go.mod services/catalog-service/go.sum* ./services/catalog-service/
COPY services/cart-service/go.mod services/cart-service/go.sum* ./services/cart-service/
COPY services/order-service/go.mod services/order-service/go.sum* ./services/order-service/
COPY services/payment-service/go.mod services/payment-service/go.sum* ./services/payment-service/
COPY services/notification-service/go.mod services/notification-service/go.sum* ./services/notification-service/

# Download dependencies
RUN go mod download

# Copy source code for compiled modules: pkg + this service
WORKDIR /app
COPY pkg/ ./pkg/
COPY services/payment-service/ ./services/payment-service/

# Build binary
WORKDIR /app/services/payment-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/payment-service ./cmd/main.go

# ==========================================
# 2. Final Runtime Stage
# ==========================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/payment-service .

EXPOSE 8085
ENTRYPOINT ["./payment-service"]