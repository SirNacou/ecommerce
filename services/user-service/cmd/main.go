package main

import (
	"context"
	"log"

	"github.com/SirNacou/ecommerce/services/user-service/internal/infrastructure/config"
	"github.com/SirNacou/ecommerce/services/user-service/internal/server"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	srv, err := server.New(ctx, *cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("server runtime error: %v", err)
	}
}
