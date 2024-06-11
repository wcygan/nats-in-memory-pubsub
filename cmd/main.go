package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {
	opts := &server.Options{}

	// Initialize new natsServer with options
	natsServer, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}

	// Start the natsServer via goroutine
	go natsServer.Start()

	// Wait for natsServer to be ready for connections
	if !natsServer.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	log.Printf("NATS server is ready")

	// Connect to natsServer
	natsClient, err := nats.Connect(natsServer.ClientURL())

	log.Printf("NATS client is connected")

	if err != nil {
		panic(err)
	}

	subject := "my-subject"

	// Subscribe to the subject
	s, err := natsClient.Subscribe(subject, func(msg *nats.Msg) {
		// Print message data
		data := string(msg.Data)
		fmt.Println(data)

		// Shutdown the natsServer (optional)
		natsServer.Shutdown()
	})

	log.Printf("Subscribed to subject: %s", subject)

	if s == nil {
		log.Fatalf("subscription is nil")
	}

	if err != nil {
		log.Fatalf("error subscribing: %v", err)
	}

	log.Printf("Unsubscribed from subject: %s", subject)

	// Publish data to the subject
	err = natsClient.Publish(subject, []byte("Hello embedded NATS!"))
	if err != nil {
		log.Fatalf("error publishing: %v", err)
	}

	// Wait for natsServer shutdown
	natsServer.WaitForShutdown()
}
