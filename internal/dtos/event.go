package dtos

import "time"

type Message interface {
	Key() string
}

type CreatedMessage struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

func (m *CreatedMessage) Key() string {
	return "message.created"
}
