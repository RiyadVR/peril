package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
		log.Fatalf("could not establish the channel: %v", err)
	}

	// _, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.Durable)
	// if err != nil {
	// 	log.Fatalf("could not subscribe to game_logs: %v", err)
	// }

	err = pubsub.SubscribeGob(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.Durable, handlerLog())
	if err != nil {
		log.Fatalf("could not subscribe to log: %v", err)
	}

	// fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamelogic.PrintServerHelp()

	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}
		firstWord := inputs[0]

		switch firstWord {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Printf("couldn't publish time: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resumes games state")
			err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})

			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "quit":
			log.Println("goodbye")
			return
		default:
			fmt.Println("unknown command")
		}

	}
}
