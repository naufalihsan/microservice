package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Product, error)
	GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error)
}

type OrderStore interface {
	Create(ctx context.Context, req *pb.CreateOrderRequest, products []*pb.Product) (*pb.Order, error)
	Get(ctx context.Context, customerId, orderId string) (*pb.Order, error)
	Update(ctx context.Context, order *pb.Order) (*pb.Order, error)
}
