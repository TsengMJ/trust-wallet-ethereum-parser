package pubsub

import (
	"context"
	"sync"
)

type Subscriber struct {
	sync.Mutex
	Handler chan *Message
	Quit    chan struct{}
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		Handler: make(chan *Message, 1024),
		Quit:    make(chan struct{}),
	}
}

func (s *Subscriber) Publish(ctx context.Context, msg *Message) {
	select {
	case <-ctx.Done():
		return
	case s.Handler <- msg:

	default:
	}
}
