package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreatePayment(context.Context, *pb.Order) (string, error) {
	// connect to payment processor
	return "", nil
}
