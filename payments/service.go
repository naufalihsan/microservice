package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/naufalihsan/msvc-payments/gateway"
	"github.com/naufalihsan/msvc-payments/processor"
)

type Service struct {
	paymentProcessor processor.PaymentProcessor
	orderGateaway    gateway.OrderGateaway
}

func NewService(paymentProcessor processor.PaymentProcessor, orderGateaway gateway.OrderGateaway) *Service {
	return &Service{paymentProcessor, orderGateaway}
}

func (s *Service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	// connect to payment processor
	paymentLink, err := s.paymentProcessor.CreatePaymentLink(order)
	if err != nil {
		return "", err
	}

	// update order with the payment link
	err = s.orderGateaway.UpdateOrderAfterPaymentLink(ctx, order.CustomerId, order.Id, paymentLink)
	if err != nil {
		return "", err
	}

	return paymentLink, nil
}
