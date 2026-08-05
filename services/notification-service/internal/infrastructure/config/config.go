package config

import "github.com/caarlos0/env/v11"

type EnvConfig struct {
	Port        string `env:"NOTIFICATION_SERVICE_PORT"`
	DatabaseURL string `env:"NOTIFICATION_DATABASE_URL"`
	JWTSecret   string `env:"JWT_SECRET"`
	NATSURL     string `env:"NATS_URL"`
}

func Load() EnvConfig {
	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	return *cfg
}
