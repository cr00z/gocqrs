package nats

import (
	"bytes"
	"encoding/gob"
	"github.com/cr00z/gocqrs/internal/domain"
	"github.com/cr00z/gocqrs/internal/dtos"
	"github.com/nats-io/nats.go"
)

type NatsEventsStore struct {
	nc              *nats.Conn
	meowCreatedSub  *nats.Subscription
	meowCreatedChan chan dtos.MeowCreatedMessage
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

	if s.meowCreatedSub != nil {
		s.meowCreatedSub.Unsubscribe()
	}
	close(s.meowCreatedChan)
}

func (s *NatsEventsStore) Publish(meow domain.Meow) error {
	mes := &dtos.MeowCreatedMessage{
		ID:        meow.ID,
		Body:      meow.Body,
		CreatedAt: meow.CreatedAt,
	}
	data, err := s.writeMessage(mes)
	if err != nil {
		return err
	}
	return s.nc.Publish(mes.Key(), data)
}

// Subscribe 1st case
func (s *NatsEventsStore) Subscribe() (<-chan dtos.MeowCreatedMessage, error) {
	var err error
	meowMsg := dtos.MeowCreatedMessage{}
	ch := make(chan *nats.Msg, 64)
	s.meowCreatedChan = make(chan dtos.MeowCreatedMessage, 64)
	s.meowCreatedSub, err = s.nc.ChanSubscribe(meowMsg.Key(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg := <-ch
			s.readMessage(msg.Data, &meowMsg)
			s.meowCreatedChan <- meowMsg
		}
	}()

	return s.meowCreatedChan, nil
}

// On 2nd case - callback
func (s *NatsEventsStore) On(f func(dtos.MeowCreatedMessage)) error {
	meowMsg := dtos.MeowCreatedMessage{}
	var err error
	s.meowCreatedSub, err = s.nc.Subscribe(
		meowMsg.Key(),
		func(msg *nats.Msg) {
			s.readMessage(msg.Data, &meowMsg)
			f(meowMsg)
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
