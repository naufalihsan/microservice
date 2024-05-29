package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/naufalihsan/msvc-common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type GrpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service OrderService
	channel *amqp.Channel
}

func NewGrpcHandler(s *grpc.Server, service OrderService, channel *amqp.Channel) {
	handler := &GrpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(s, handler)
}

func (h *GrpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("new order received 🛒 from Customer", req.CustomerId)

	if err := h.service.ValidateOrder(ctx, req); err != nil {
		return nil, err
	}

	order := &pb.Order{
		Id:         "1",
		CustomerId: req.CustomerId,
	}

	jsonOrder, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	queue, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	h.channel.PublishWithContext(ctx, "", queue.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         jsonOrder,
		DeliveryMode: amqp.Persistent,
	})

	return order, nil
}
