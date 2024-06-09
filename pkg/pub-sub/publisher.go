package pubsub

import (
	"context"
	"sync"
)

type Publisher struct {
	sync.Mutex
	subs map[*Subscriber]struct{}
}

type Message struct {
	// Data []byte
	Data interface{}
}

var DefaultPublisher *Publisher = NewPublisher()

func NewPublisher() *Publisher {
	return &Publisher{
		subs: map[*Subscriber]struct{}{},
	}
}

func (p *Publisher) Subscribe(ctx context.Context, s *Subscriber) error {
	p.Lock()
	p.subs[s] = struct{}{}
	p.Unlock()

	return nil
}

func (p *Publisher) Unsubscribe(ctx context.Context, s *Subscriber) error {
	p.Lock()
	delete(p.subs, s)
	p.Unlock()

	return nil
}

func (p *Publisher) Publish(ctx context.Context, msg *Message) error {
	p.Lock()

	for s := range p.subs {
		s.Publish(ctx, msg)
	}
	p.Unlock()

	return nil
}
