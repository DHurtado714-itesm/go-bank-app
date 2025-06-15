package messagebus

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	natsdriver "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
)

// NATSBus implements MessageBus using NATS as backend via Watermill.
type NATSBus struct {
	publisher  *natsdriver.Publisher
	subscriber *natsdriver.Subscriber
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewNATSBus creates a new NATSBus connected to the provided URL.
func NewNATSBus(url string) (*NATSBus, error) {
	logger := watermill.NewStdLogger(false, false)

	marshaler := &natsdriver.JSONMarshaler{}

	pub, err := natsdriver.NewPublisher(natsdriver.PublisherConfig{
		URL:       url,
		Marshaler: marshaler,
	}, logger)
	if err != nil {
		return nil, err
	}

	sub, err := natsdriver.NewSubscriber(natsdriver.SubscriberConfig{
		URL:         url,
		Unmarshaler: marshaler,
	}, logger)
	if err != nil {
		_ = pub.Close()
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &NATSBus{publisher: pub, subscriber: sub, ctx: ctx, cancel: cancel}, nil
}

// Close gracefully closes the underlying connections.
func (b *NATSBus) Close() error {
	b.cancel()
	if err := b.publisher.Close(); err != nil {
		return err
	}
	return b.subscriber.Close()
}

// Publish publishes data as JSON on the given topic.
func (b *NATSBus) Publish(topic string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), bytes)
	return b.publisher.Publish(topic, msg)
}

// Subscribe registers handler for the topic. Handler will be called in a new goroutine for each message.
func (b *NATSBus) Subscribe(topic string, handler func(msg *message.Message)) error {
	messages, err := b.subscriber.Subscribe(b.ctx, topic)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			handler(msg)
			msg.Ack()
		}
	}()

	return nil
}
