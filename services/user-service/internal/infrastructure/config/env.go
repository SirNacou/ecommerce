package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type EnvConfig struct {
	DatabaseURL       string        `env:"DATABASE_URL,required"`
	Port              string        `env:"PORT,required"`
	JWTSecret         string        `env:"JWT_SECRET,required"`
	AccessExpiration  time.Duration `env:"ACCESS_EXPIRATION" envDefault:"15m"`
	RefreshExpiration time.Duration `env:"REFRESH_EXPIRATION" envDefault:"24h"`
}

func LoadEnvConfig() (*EnvConfig, error) {
	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
