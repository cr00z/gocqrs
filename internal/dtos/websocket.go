package dtos

import "time"

const KindMessageCreated = iota + 1

type CreatedMessageWs struct {
	Kind      uint32    `json:"kind"`
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCreatedMessageWs(id string, body string, createdAt time.Time) *CreatedMessageWs {
	return &CreatedMessageWs{
		Kind:      KindMessageCreated,
		ID:        id,
		Body:      body,
		CreatedAt: createdAt,
	}
}
