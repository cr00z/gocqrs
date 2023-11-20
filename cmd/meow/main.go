package main

import (
	"github.com/cr00z/gocqrs/cmd/meow/config"
	"github.com/cr00z/gocqrs/internal/controller/nats"
	"github.com/cr00z/gocqrs/internal/repository/postgres"
	"github.com/cr00z/gocqrs/pkg/util"
)

func main() {
	cfg := util.Must(config.NewMeowConfig)

	pgRepo := util.Must(func() (*postgres.PostgresRepository, error) {
		return postgres.NewPostgresRepository(cfg.PostgresDsn)
	})
	defer pgRepo.Close()

	natsCtrl := util.Must(func() (*nats.NatsEventsStore, error) {
		return nats.NewNats(cfg.NatsUrl)
	})
	defer natsCtrl.Close()
}
