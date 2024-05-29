package gateway

import (
	"context"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
	discovery "github.com/naufalihsan/msvc-common/discovery"
)

type GrpcGateway struct {
	registry discovery.Registry
}

func NewGrpcGateway(registry discovery.Registry) *GrpcGateway {
	return &GrpcGateway{registry}
}

func (g *GrpcGateway) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, common.OrdersService, g.registry)
	if err != nil {
		return nil, err
	}

	client := pb.NewOrderServiceClient(conn)

	return client.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId:    req.CustomerId,
		OrderProducts: req.OrderProducts,
	})
}
