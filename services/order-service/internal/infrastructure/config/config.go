package config

import (
	"github.com/caarlos0/env/v11"
)

type EnvConfig struct {
	Port        string `env:"ORDER_SERVICE_PORT"`
	DatabaseURL string `env:"ORDER_DATABASE_URL"`
	JWTSecret   string `env:"JWT_SECRET"`
	NATSURL     string `env:"NATS_URL"`
	CatalogURL  string `env:"CATALOG_SERVICE_URL"`
}

func Load() (*EnvConfig, error) {
	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
