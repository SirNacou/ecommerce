package main

import (
	"context"
	"log"

	"github.com/SirNacou/ecommerce/services/order-service/internal/infrastructure/config"
	"github.com/SirNacou/ecommerce/services/order-service/internal/infrastructure/server"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	srv, err := server.New(ctx, *cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
