package producer

import (
	"fmt"

	"github.com/streadway/amqp"
)

type MessageProducer interface {
	Publish([]byte) error
	Close() error
}

type RabbitProducer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func NewRabbitProducer(url, queueName string) (*RabbitProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &RabbitProducer{
		conn:  conn,
		ch:    ch,
		queue: queueName,
	}, nil
}

func (r *RabbitProducer) Publish(msg []byte) error {
	return r.ch.Publish(
		"",
		r.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)
}

func (r *RabbitProducer) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}
