package main

import (
	"log"

	"github.com/cr00z/gocqrs/internal/controller/http/createhttp"
	"github.com/cr00z/gocqrs/internal/controller/nats"
	"github.com/cr00z/gocqrs/internal/repository/postgres"
	"github.com/cr00z/gocqrs/pkg/config"
	"github.com/cr00z/gocqrs/pkg/util"
)

func main() {
	cfg := util.Must(config.NewConfig)

	pgRepo := util.MustStr(postgres.NewPostgresRepository, cfg.PostgresDsn)
	defer pgRepo.Close()

	natsCtrl := util.MustStr(nats.NewNats, cfg.NatsUrl)
	defer natsCtrl.Close()

	err := createhttp.NewHttpServer(pgRepo, natsCtrl).Start()
	if err != nil {
		log.Fatal(err)
	}
}
