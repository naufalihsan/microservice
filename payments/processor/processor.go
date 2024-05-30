package processor

import pb "github.com/naufalihsan/msvc-common/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
