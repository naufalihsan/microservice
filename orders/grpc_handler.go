package main

import (
	"context"
	"log"

	pb "github.com/naufalihsan/msvc-common/api"
	"google.golang.org/grpc"
)

type GrpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGrpcHandler(s *grpc.Server) {
	handler := &GrpcHandler{}
	pb.RegisterOrderServiceServer(s, handler)
}

func (h *GrpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("new order received 🛒 from Customer", req.CustomerId)

	order := &pb.Order{
		Id:         "1",
		CustomerId: req.CustomerId,
	}

	return order, nil
}
