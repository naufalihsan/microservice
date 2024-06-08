package main

import (
	"context"
	"fmt"

	pb "github.com/naufalihsan/msvc-common/api"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	next InventoryService
}

func NewTelemetry(next InventoryService) *Telemetry {
	return &Telemetry{next}
}

func (t *Telemetry) Get(ctx context.Context, ids []string) ([]*pb.Product, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("Get Inventory: %v", ids))

	return t.next.Get(ctx, ids)
}

func (t *Telemetry) Validate(ctx context.Context, req []*pb.OrderProduct) (bool, []*pb.Product, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("Validate Inventory: %v", req))

	return t.next.Validate(ctx, req)
}

func (t *Telemetry) Update(ctx context.Context, products []*pb.Product) error {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("Update Inventory: %v", products))

	return t.next.Update(ctx, products)
}
