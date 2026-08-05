FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod* go.sum* go.work* ./
COPY pkg/go.mod* pkg/go.sum* ./pkg/
COPY services/notification-service/go.mod services/notification-service/go.sum* ./services/notification-service/

WORKDIR /app/services/notification-service
RUN GOWORK=off go mod download

WORKDIR /app
COPY pkg/ ./pkg/
COPY services/notification-service/ ./services/notification-service/

WORKDIR /app/services/notification-service
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -ldflags="-w -s" -o /app/notification-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/notification-service .

EXPOSE 8086
ENTRYPOINT ["./notification-service"]