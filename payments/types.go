package main

import (
	"context"

	pb "github.com/naufalihsan/msvc-common/api"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
