package main

import (
	"context"
	"log"

	"github.com/SirNacou/ecommerce/services/notification-service/internal/infrastructure/config"
	"github.com/SirNacou/ecommerce/services/notification-service/internal/infrastructure/server"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	srv, err := server.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
