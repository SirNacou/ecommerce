FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod* go.sum* go.work* ./
COPY pkg/go.mod* pkg/go.sum* ./pkg/
COPY services/cart-service/go.mod services/cart-service/go.sum* ./services/cart-service/

WORKDIR /app/services/cart-service
RUN GOWORK=off go mod download

WORKDIR /app
COPY pkg/ ./pkg/
COPY services/cart-service/ ./services/cart-service/

WORKDIR /app/services/cart-service
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -ldflags="-w -s" -o /app/cart-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/cart-service .

EXPOSE 8083
ENTRYPOINT ["./cart-service"]