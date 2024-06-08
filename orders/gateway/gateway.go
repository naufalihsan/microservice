package gateway

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type InventoryGateaway interface {
	Validate(ctx context.Context, orderProducts []*pb.OrderProduct) ([]*pb.Product, error)
}
