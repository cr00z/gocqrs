package controller

import (
	"github.com/cr00z/gocqrs/internal/domain"
	"github.com/cr00z/gocqrs/internal/dtos"
)

type EventStore interface {
	Close()
	Publish(meow domain.Meow) error
	Subscribe(<-chan dtos.MeowCreatedMessage, error)
	On(f func(dtos.MeowCreatedMessage)) error
}
