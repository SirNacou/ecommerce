package config

import (
	"github.com/caarlos0/env/v11"
)

type EnvConfig struct {
	Port        string `env:"PORT,required"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	JWTSecret   string `env:"JWT_SECRET,required"`
}

func Load() (*EnvConfig, error) {
	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
