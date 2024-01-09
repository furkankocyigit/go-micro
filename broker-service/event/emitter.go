package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel,err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}	

func (e *Emitter) Push(event string, severity string) error {
	channel,err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Printf("Pushing event %s with severity %s", event, severity)

	err = channel.Publish(
		"logs_topic", // exchange
		severity,     // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		})
	if err != nil {
		return err
	}

	return nil
}

func NewEmitter(conn *amqp.Connection) (Emitter,error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}
	
	return emitter, nil
}

	