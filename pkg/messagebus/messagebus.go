package messagebus

import "github.com/ThreeDotsLabs/watermill/message"

// MessageBus is a generic interface for publishing and subscribing to messages.
type MessageBus interface {
	// Publish sends the given data to the specified topic.
	Publish(topic string, data interface{}) error

	// Subscribe registers a handler for messages on the given topic.
	Subscribe(topic string, handler func(msg *message.Message)) error
}
