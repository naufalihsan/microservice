package main

import (
	"context"
	"testing"

	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/naufalihsan/msvc-payments/processor/inmem"
)

func TestService(t *testing.T) {
	inmemProcessor := inmem.NewProcessor()
	service := NewService(inmemProcessor)

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
