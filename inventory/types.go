package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type InventoryService interface {
	Get(ctx context.Context, ids []string) ([]*pb.Product, error)
	Validate(ctx context.Context, orderProducts []*pb.OrderProduct) ([]*pb.Product, error)
	Update(ctx context.Context, products []*pb.Product) error
}

type InventoryStore interface {
	Get(ctx context.Context, ids []string) ([]*pb.Product, error)
	Update(ctx context.Context, inventoryMap map[string]*pb.Product) error
}
