package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/naufalihsan/msvc-payments/processor"
)

type Service struct {
	paymentProcessor processor.PaymentProcessor
}

func NewService(paymentProcessor processor.PaymentProcessor) *Service {
	return &Service{paymentProcessor}
}

func (s *Service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	// connect to payment processor
	return s.paymentProcessor.CreatePaymentLink(order)
}
