package discovery

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addresses, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(
		addresses[rand.Intn(len(addresses))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
