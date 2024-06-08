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

func (g *GrpcGateway) Validate(ctx context.Context, orderProducts []*pb.OrderProduct) ([]*pb.Product, error) {
	conn, err := discovery.ServiceConnection(ctx, common.InventoryService, g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewInventoryServiceClient(conn)

	res, err := client.Validate(ctx, &pb.ValidateInventoryRequest{
		OrderProducts: orderProducts,
	})

	if err != nil {
		return nil, err
	}

	return res.Products, nil
}
