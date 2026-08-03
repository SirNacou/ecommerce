# ==========================================
# 1. Build Stage
# ==========================================
FROM golang:1.26-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o user-service ./cmd/main.go

# ==========================================
# 2. Final Runtime Stage
# ==========================================
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/user-service .
EXPOSE 8080
ENTRYPOINT ["./user-service"]