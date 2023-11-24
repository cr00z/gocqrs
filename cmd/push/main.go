package main

import (
	"log"

	"github.com/cr00z/gocqrs/internal/controller/http/pushhttp"
	"github.com/cr00z/gocqrs/internal/controller/nats"
	"github.com/cr00z/gocqrs/internal/controller/websock"
	"github.com/cr00z/gocqrs/internal/dtos"
	"github.com/cr00z/gocqrs/pkg/config"
	"github.com/cr00z/gocqrs/pkg/util"
)

func main() {
	cfg := util.Must(config.NewConfig)

	hub := websock.NewHub()
	go hub.Run()

	log.Printf("Nats: %s", cfg.NatsUrl)
	natsCtrl := util.MustStr(nats.NewNats, cfg.NatsUrl)
	defer natsCtrl.Close()
	util.Must(func() (any, error) {
		onMsgCreated := func(m dtos.CreatedMessage) {
			log.Printf("message received: %v", m)
			hub.Broadcast(dtos.NewCreatedMessageWs(m.ID, m.Body, m.CreatedAt), nil)
		}
		return nil, natsCtrl.On(onMsgCreated)
	})

	err := pushhttp.NewHttpServer(hub).Start()
	if err != nil {
		log.Fatal(err)
	}
}
