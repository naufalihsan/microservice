package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, order *pb.Order) (string, error)
}
