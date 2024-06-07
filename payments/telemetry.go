package main

import (
	"context"
	"fmt"

	pb "github.com/naufalihsan/msvc-common/api"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	next PaymentService
}

func NewTelemetry(next PaymentService) *Telemetry {
	return &Telemetry{next}
}

func (t *Telemetry) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("Create Payment: %v", order))

	return t.next.CreatePayment(ctx, order)
}
