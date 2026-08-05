FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod* go.sum* go.work* ./
COPY pkg/go.mod* pkg/go.sum* ./pkg/
COPY services/order-service/go.mod services/order-service/go.sum* ./services/order-service/

WORKDIR /app/services/order-service
RUN GOWORK=off go mod download

WORKDIR /app
COPY pkg/ ./pkg/
COPY services/order-service/ ./services/order-service/

WORKDIR /app/services/order-service
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -ldflags="-w -s" -o /app/order-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/order-service .

EXPOSE 8084
ENTRYPOINT ["./order-service"]