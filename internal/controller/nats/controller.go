package nats

import (
	"bytes"
	"encoding/gob"

	"github.com/cr00z/gocqrs/internal/domain"
	"github.com/cr00z/gocqrs/internal/dtos"
	"github.com/nats-io/nats.go"
)

type NatsEventsStore struct {
	nc          *nats.Conn
	createdSub  *nats.Subscription
	createdChan chan dtos.CreatedMessage
}

func NewNats(url string) (*NatsEventsStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventsStore{
		nc: nc,
	}, nil
}

func (s *NatsEventsStore) Close() {
	if s.nc != nil {
		s.nc.Close()
	}

	if s.createdSub != nil {
		_ = s.createdSub.Unsubscribe()
	}
	close(s.createdChan)
}

func (s *NatsEventsStore) Publish(message domain.Message) error {
	mes := &dtos.CreatedMessage{
		ID:        message.ID,
		Body:      message.Body,
		CreatedAt: message.CreatedAt,
	}
	data, err := s.writeMessage(mes)
	if err != nil {
		return err
	}
	return s.nc.Publish(mes.Key(), data)
}

// Subscribe 1st case
func (s *NatsEventsStore) Subscribe() (<-chan dtos.CreatedMessage, error) {
	var err error
	newMsg := dtos.CreatedMessage{}
	ch := make(chan *nats.Msg, 64)
	s.createdChan = make(chan dtos.CreatedMessage, 64)
	s.createdSub, err = s.nc.ChanSubscribe(newMsg.Key(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg := <-ch
			_ = s.readMessage(msg.Data, &newMsg)
			s.createdChan <- newMsg
		}
	}()

	return s.createdChan, nil
}

// On 2nd case - callback
func (s *NatsEventsStore) On(f func(dtos.CreatedMessage)) error {
	newMsg := dtos.CreatedMessage{}
	var err error
	s.createdSub, err = s.nc.Subscribe(
		newMsg.Key(),
		func(msg *nats.Msg) {
			_ = s.readMessage(msg.Data, &newMsg)
			f(newMsg)
		})
	return err
}

func (s *NatsEventsStore) readMessage(data []byte, m any) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

func (s *NatsEventsStore) writeMessage(m dtos.Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
