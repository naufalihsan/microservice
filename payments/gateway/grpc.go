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

func (g *GrpcGateway) UpdateOrderAfterPaymentLink(ctx context.Context, customerId, orderId, paymentLink string) error {
	conn, err := discovery.ServiceConnection(ctx, common.OrderService, g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	_, err = client.UpdateOrder(ctx, &pb.Order{
		Id:          orderId,
		CustomerId:  customerId,
		Status:      common.OrderStatusWaitingPayment,
		PaymentLink: paymentLink,
	})

	return err
}
