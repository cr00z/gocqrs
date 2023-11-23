package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

type Config struct {
	PostgresDsn string
	NatsUrl     string
	ElasticUrl  string
}

func NewConfig() (*Config, error) {
	var config config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &Config{
		PostgresDsn: fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable",
			config.PostgresUser, config.PostgresPassword, config.PostgresDB),
		NatsUrl:    fmt.Sprintf("nats://%s", config.NatsAddress),
		ElasticUrl: fmt.Sprintf("https://%s", config.ElasticsearchAddress),
	}, nil
}
