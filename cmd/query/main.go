package main

import (
	"context"
	"log"

	"github.com/cr00z/gocqrs/internal/controller/http/createhttp"
	"github.com/cr00z/gocqrs/internal/controller/nats"
	"github.com/cr00z/gocqrs/internal/domain"
	"github.com/cr00z/gocqrs/internal/dtos"
	"github.com/cr00z/gocqrs/internal/repository/elastic"
	"github.com/cr00z/gocqrs/internal/repository/postgres"
	"github.com/cr00z/gocqrs/pkg/config"
	"github.com/cr00z/gocqrs/pkg/util"
)

func main() {
	ctx := context.Background()

	cfg := util.Must(config.NewConfig)

	log.Printf("Postgres: %s", cfg.PostgresDsn)
	pgRepo := util.MustStr(postgres.NewPostgresRepository, cfg.PostgresDsn)
	defer pgRepo.Close()

	log.Printf("Elastic: %s", cfg.ElasticUrl)
	elRepo := util.MustStr(elastic.NewElasticRepository, cfg.ElasticUrl)
	defer elRepo.Close()

	log.Printf("Nats: %s", cfg.NatsUrl)
	natsCtrl := util.MustStr(nats.NewNats, cfg.NatsUrl)
	defer natsCtrl.Close()
	util.Must(func() (any, error) {
		onMsgCreated := func(m dtos.CreatedMessage) {
			msg := domain.Message{
				ID:        m.ID,
				Body:      m.Body,
				CreatedAt: m.CreatedAt,
			}
			if err := elRepo.Insert(ctx, msg); err != nil {
				log.Print(err)
			}
		}
		return nil, natsCtrl.On(onMsgCreated)
	})

	err := createhttp.NewHttpServer(pgRepo, natsCtrl).Start()
	if err != nil {
		log.Fatal(err)
	}
}
