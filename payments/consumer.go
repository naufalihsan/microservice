package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/naufalihsan/msvc-common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	service PaymentService
}

func NewConsumer(service PaymentService) *Consumer {
	return &Consumer{service}
}

func (c *Consumer) Listen(channel *amqp.Channel) {
	queue, err := channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			log.Printf("message received %v", message)

			order := &pb.Order{}
			if err := json.Unmarshal(message.Body, order); err != nil {
				log.Fatal(err)
			}

			c.service.CreatePayment(context.Background(), order)
		}
	}()

	<-forever

}
