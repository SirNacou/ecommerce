package main

import (
	"context"
	"log"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/infrastructure/config"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/infrastructure/server"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	srv, err := server.New(ctx, *cfg)
	if err != nil {
		log.Fatalf("failed to initialize catalog-service: %v", err)
	}

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("catalog-service runtime error: %v", err)
	}
}
