package inmem

import (
	"fmt"

	pb "github.com/naufalihsan/msvc-common/api"
)

type InMem struct{}

func NewProcessor() *InMem {
	return &InMem{}
}

func (i *InMem) CreatePaymentLink(order *pb.Order) (string, error) {
	return fmt.Sprintf("local-%s", order.CustomerId), nil
}
