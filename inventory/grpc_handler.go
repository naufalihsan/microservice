package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type GrpcHandler struct {
	pb.UnimplementedInventoryServiceServer

	service InventoryService
	channel *amqp.Channel
}

func NewGrpcHandler(s *grpc.Server, service InventoryService, channel *amqp.Channel) {
	handler := &GrpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterInventoryServiceServer(s, handler)
}

func (h *GrpcHandler) Get(ctx context.Context, req *pb.GetInventoryRequest) (*pb.GetInventoryResponse, error) {
	products, err := h.service.Get(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	return &pb.GetInventoryResponse{Products: products}, nil
}

func (h *GrpcHandler) Validate(ctx context.Context, req *pb.ValidateInventoryRequest) (*pb.ValidateInventoryResponse, error) {
	products, err := h.service.Validate(ctx, req.OrderProducts)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateInventoryResponse{Status: true, Products: products}, nil
}
