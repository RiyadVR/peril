package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queuename,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {

	ch, queue, err := DeclareAndBind(conn, exchange, queuename, key, queueType)
	if err != nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			var data T
			err := json.Unmarshal(msg.Body, &data)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}

			acktype := handler(data)
			switch acktype {
			case Ack:
				msg.Ack(false)
				fmt.Println("Acknowledged")
			case NackRequeue:
				msg.Nack(false, true)
				fmt.Println("NackRequeued")
			case NackDiscard:
				msg.Nack(false, false)
				fmt.Println("NackDiscarded")
			default:
				fmt.Printf("unknown Acktype: %v", acktype)
			}

		}
	}()

	return nil
}
