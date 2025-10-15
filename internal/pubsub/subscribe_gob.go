package pubsub

import (
	"bytes"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](
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

	err = ch.Qos(10, 0, false)
	if err != nil {
		return fmt.Errorf("could not set QoS: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			var data T
			buff := bytes.NewBuffer(msg.Body)
			dec := gob.NewDecoder(buff)
			err := dec.Decode(&data)
			if err != nil {
				fmt.Printf("could not decode gob message: %v\n", err)
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
