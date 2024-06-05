package main

import (
	"context"
	"testing"

	pb "github.com/naufalihsan/msvc-common/api"
	inmemRegistry "github.com/naufalihsan/msvc-common/discovery/inmem"
	"github.com/naufalihsan/msvc-payments/gateway"
	inmemProcessor "github.com/naufalihsan/msvc-payments/processor/inmem"
)

func TestService(t *testing.T) {
	inmemProcessor := inmemProcessor.NewProcessor()
	inmemRegistry, _ := inmemRegistry.NewRegistry()
	orderGateaway := gateway.NewGrpcGateway(inmemRegistry)

	service := NewService(inmemProcessor, orderGateaway)

	t.Run("should create payment link", func(t *testing.T) {
		paymentLink, err := service.CreatePayment(context.Background(), &pb.Order{})

		if err != nil {
			t.Errorf("Create Payment Error %v", err)
		}

		if paymentLink == "" {
			t.Error("Create payment link is empty")
		}
	})
}
