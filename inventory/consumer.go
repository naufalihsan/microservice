package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/naufalihsan/msvc-common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type Consumer struct {
	service InventoryService
}

func NewConsumer(service InventoryService) *Consumer {
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
			// extract the headers
			amqpCtx := broker.ExtractHeader(context.Background(), message.Headers)
			tracer := otel.Tracer("amqp")
			_, span := tracer.Start(amqpCtx, fmt.Sprintf("AMQP Consume %s", queue.Name))

			log.Printf("message received %v", message)

			order := &pb.Order{}
			if err := json.Unmarshal(message.Body, order); err != nil {
				message.Nack(false, false)
				log.Fatal(err)
			}

			err := c.service.Update(context.Background(), order.Products)
			if err != nil {
				// update status order to out of stock
				log.Printf("failed to update inventory: %v", err)
			}

			span.End()
			message.Ack(false)
		}
	}()

	<-forever
}
