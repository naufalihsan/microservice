package gateway

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type OrderGateaway interface {
	UpdateOrderAfterPaid(ctx context.Context, order *pb.Order) error
}
