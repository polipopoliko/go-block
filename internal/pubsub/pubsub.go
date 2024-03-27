package pubsub

import (
	"context"

	ps "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type TopicHandler interface {
	TopicPublisher(ctx context.Context) error
}

type PubSub struct {
	topic string
	host  host.Host
}

func NewPubSub(h host.Host, topic string) *PubSub {
	return &PubSub{
		topic: topic,
		host:  h,
	}
}

func (p *PubSub) GetTopic(ctx context.Context) (*ps.Topic, error) {
	ps, err := ps.NewGossipSub(ctx, p.host)
	if err != nil {
		return nil, err
	}

	topic, err := ps.Join(p.topic)
	if err != nil {
		return nil, err
	}

	return topic, err
}
