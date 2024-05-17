package main

import (
	"context"
	"log"

	pb "github.com/naufalihsan/msvc-common/api"
	"google.golang.org/grpc"
)

type GrpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
}

func NewGrpcHandler(s *grpc.Server, service OrderService) {
	handler := &GrpcHandler{
		service: service,
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

	return order, nil
}
