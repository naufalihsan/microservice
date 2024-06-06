package broker

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MaxRetryCount = 3
)

func Connect(user, pass, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
	conn, err := amqp.Dial(address)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = channel.ExchangeDeclare(OrderCreatedEvent, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = channel.ExchangeDeclare(OrderPaidEvent, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = channel.ExchangeDeclare(DeadLetterExchange, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	dlq, err := channel.QueueDeclare(DeadLetterQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = channel.QueueBind(dlq.Name, "", DeadLetterExchange, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return channel, conn.Close
}

func HandleRetry(channel *amqp.Channel, message *amqp.Delivery) error {
	if message.Headers == nil {
		message.Headers = amqp.Table{}
	}

	retryCount, _ := message.Headers["x-retry-count"].(int64)
	retryCount++
	message.Headers["x-retry-count"] = retryCount

	log.Printf("Retrying message: %s, attemp: %d", message.Body, retryCount)

	if retryCount >= MaxRetryCount {
		return channel.PublishWithContext(context.Background(), "", DeadLetterQueue, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Headers:      message.Headers,
			Body:         message.Body,
			DeliveryMode: amqp.Persistent,
		})
	}

	time.Sleep(time.Second * time.Duration(retryCount))

	return channel.PublishWithContext(context.Background(), message.Exchange, message.RoutingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Headers:      message.Headers,
		Body:         message.Body,
		DeliveryMode: amqp.Persistent,
	})
}
