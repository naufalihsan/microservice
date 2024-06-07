package main

import (
	"context"
	"fmt"

	pb "github.com/naufalihsan/msvc-common/api"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	next OrderService
}

func NewTelemetry(next OrderService) *Telemetry {
	return &Telemetry{next}
}

func (t *Telemetry) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("Create Order: %v", req))

	return t.next.CreateOrder(ctx, req)
}

func (t *Telemetry) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Product, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("Validate Order: %v", req))

	return t.next.ValidateOrder(ctx, req)
}

func (t *Telemetry) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("Get Order: %v", req))

	return t.next.GetOrder(ctx, req)
}

func (t *Telemetry) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("Update Order: %v", order))

	return t.next.UpdateOrder(ctx, order)
}
