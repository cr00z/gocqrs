package controller

import (
	"github.com/cr00z/gocqrs/internal/domain"
	"github.com/cr00z/gocqrs/internal/dtos"
)

type EventStore interface {
	Close()
	Publish(message domain.Message) error
	Subscribe() (<-chan dtos.CreatedMessage, error)
	On(f func(dtos.CreatedMessage)) error
}
