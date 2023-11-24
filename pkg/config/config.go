package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	PostgresHost         string `envconfig:"POSTGRES_HOST"         default:"localhost:5432"`
	PostgresDB           string `envconfig:"POSTGRES_DB"           default:"messages"`
	PostgresUser         string `envconfig:"POSTGRES_USER"         default:"messages"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"     default:"password"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"          default:"localhost:4222"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS" default:"localhost:9200"`
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
		PostgresDsn: fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresDB),
		NatsUrl:    fmt.Sprintf("nats://%s", config.NatsAddress),
		ElasticUrl: fmt.Sprintf("http://%s", config.ElasticsearchAddress),
	}, nil
}
