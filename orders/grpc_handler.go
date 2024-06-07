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
	queue, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	tracer := otel.Tracer("amqp")
	amqpCtx, span := tracer.Start(ctx, fmt.Sprintf("AMQP Publish %s", queue.Name))
	defer span.End()

	order, err := h.service.CreateOrder(amqpCtx, req)
	if err != nil {
		return nil, err
	}

	jsonOrder, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	// inject the headers
	headers := broker.InjectHeaders(amqpCtx)

	h.channel.PublishWithContext(amqpCtx, "", queue.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         jsonOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
	})

	return order, nil
}

func (h *GrpcHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, req)
}

func (h *GrpcHandler) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, order)
}
