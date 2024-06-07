package main

import (
	"context"
	"time"

	pb "github.com/naufalihsan/msvc-common/api"
	"go.uber.org/zap"
)

type Logging struct {
	next OrderService
}

func NewLogging(next OrderService) *Logging {
	return &Logging{next}
}

func (t *Logging) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CreateOrder", zap.Duration("took", time.Since(start)))
	}()

	return t.next.CreateOrder(ctx, req)
}

func (t *Logging) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Product, error) {
	return t.next.ValidateOrder(ctx, req)
}

func (t *Logging) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return t.next.GetOrder(ctx, req)
}

func (t *Logging) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	return t.next.UpdateOrder(ctx, order)
}
