package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill/message"
	"go-bank-app/pkg/messagebus"
)

type UserCreated struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func main() {
	bus, err := messagebus.NewNATSBus("nats://localhost:4222")
	if err != nil {
		log.Fatalf("failed to create message bus: %v", err)
	}
	defer bus.Close()

	// Subscriber
	if err := bus.Subscribe("users.created", func(msg *message.Message) {
		var event UserCreated
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			log.Printf("could not unmarshal event: %v", err)
			return
		}
		log.Printf("received user: %#v", event)
	}); err != nil {
		log.Fatalf("subscribe error: %v", err)
	}

	// Publisher example
	go func() {
		event := UserCreated{ID: "123", Email: "test@example.com"}
		if err := bus.Publish("users.created", event); err != nil {
			log.Printf("publish error: %v", err)
		}
	}()

	// Wait for termination signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	// Give some time for graceful shutdown
	_ = bus.Close()
}
