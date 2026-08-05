package config

import (
	"github.com/caarlos0/env/v11"
)

type EnvConfig struct {
	Port        string `env:"CART_SERVICE_PORT"`
	DatabaseURL string `env:"CART_DATABASE_URL"`
	JWTSecret   string `env:"JWT_SECRET"`
	CatalogURL  string `env:"CATALOG_SERVICE_URL"`
}

func Load() (*EnvConfig, error) {
	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil

}
