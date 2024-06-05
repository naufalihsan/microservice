package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

var (
	OrderStatusPending = "pending"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Product, error)
	GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error)
}

type OrderStore interface {
	Create(ctx context.Context, req *pb.CreateOrderRequest, products []*pb.Product) (*pb.Order, error)
	Get(ctx context.Context, customerId, orderId string) (*pb.Order, error)
}
