package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Product, error)
}

type OrderStore interface {
	Create(ctx context.Context) error
}
