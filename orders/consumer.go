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
	service OrderService
}

func NewConsumer(service OrderService) *Consumer {
	return &Consumer{service}
}

func (c *Consumer) Listen(channel *amqp.Channel, instanceId string) {
	// Declares queue
	queue, err := channel.QueueDeclare(instanceId, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Binds this queue to an exchange (fanout)
	err = channel.QueueBind(queue.Name, "", broker.OrderPaidEvent, false, nil)
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
				message.Nack(false, false)
				log.Fatal(err)
			}

			order, err := c.service.UpdateOrder(context.Background(), order)
			if err != nil {
				log.Printf("failed to create payment: %v", err)

				if err := broker.HandleRetry(channel, &message); err != nil {
					log.Printf("error handling retry: %v", err)
				}

				continue
			}

			message.Ack(false)

			log.Printf("🐇 order %s has been updated", order.Id)
		}
	}()

	<-forever
}
