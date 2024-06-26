package gateway

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type OrderGateaway interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(ctx context.Context, customerId, orderId string) (*pb.Order, error)
}
