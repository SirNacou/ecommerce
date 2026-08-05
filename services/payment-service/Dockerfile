FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod* go.sum* go.work* ./
COPY pkg/go.mod* pkg/go.sum* ./pkg/
COPY services/payment-service/go.mod services/payment-service/go.sum* ./services/payment-service/

WORKDIR /app/services/payment-service
RUN GOWORK=off go mod download

WORKDIR /app
COPY pkg/ ./pkg/
COPY services/payment-service/ ./services/payment-service/

WORKDIR /app/services/payment-service
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -ldflags="-w -s" -o /app/payment-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/payment-service .

EXPOSE 8085
ENTRYPOINT ["./payment-service"]