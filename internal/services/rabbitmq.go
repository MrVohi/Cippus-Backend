package services

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitPublisher(url string) (*RabbitPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	if _, err = ch.QueueDeclare("notifications.dispatch", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	return &RabbitPublisher{conn: conn, channel: ch}, nil
}

func (r *RabbitPublisher) Publish(queueName string, payload []byte) error {
	return r.channel.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}

func (r *RabbitPublisher) Close() {
	r.channel.Close()
	r.conn.Close()
}
