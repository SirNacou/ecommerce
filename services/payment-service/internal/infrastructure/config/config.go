package config

import "github.com/caarlos0/env/v11"

type EnvConfig struct {
	Port        string `env:"PAYMENT_SERVICE_PORT"`
	DatabaseURL string `env:"PAYMENT_DATABASE_URL"`
	JWTSecret   string `env:"JWT_SECRET"`
	NATSURL     string `env:"NATS_URL"`
}

func Load() (*EnvConfig, error) {
	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}