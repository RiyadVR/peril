package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to the RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection to the message broker is successful")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not establisht the channel: %v", err)
	}

	err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})

	if err != nil {
		log.Printf("could not publish time: %v", err)
	}

	fmt.Println("Pause message sent!")

}
