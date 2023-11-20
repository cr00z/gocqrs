package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

type meowConfig struct {
	PostgresDsn string
	NatsUrl     string
}

func NewMeowConfig() (*meowConfig, error) {
	var config config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &meowConfig{
		PostgresDsn: fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable",
			config.PostgresUser, config.PostgresPassword, config.PostgresDB),
		NatsUrl: fmt.Sprintf("nats://%s", config.NatsAddress),
	}, nil
}
