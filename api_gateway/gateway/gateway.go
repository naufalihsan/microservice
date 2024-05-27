package gateway

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type OrdersGateaway interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
}
