package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
)

// temp inmem store
var inmemStore = make(map[string]*pb.Order)

type Store struct {
	// connect to mongodb
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Create(ctx context.Context, req *pb.CreateOrderRequest, products []*pb.Product) (*pb.Order, error) {
	uniqueId := strconv.Itoa(rand.Intn(1000))
	customerId := req.CustomerId

	uniqueKey := GenerateUniqueKey(customerId, uniqueId)

	order := &pb.Order{
		Id:         uniqueId,
		CustomerId: customerId,
		Status:     OrderStatusPending,
		Products:   products,
	}

	inmemStore[uniqueKey] = order

	return order, nil
}

func (s *Store) Get(ctx context.Context, customerId, orderId string) (*pb.Order, error) {
	uniqueKey := GenerateUniqueKey(customerId, orderId)

	order, ok := inmemStore[uniqueKey]
	if !ok {
		return nil, common.ErrInvalidOrder
	}

	return order, nil
}

func GenerateUniqueKey(customerId, orderId string) string {
	return fmt.Sprintf("%s-%s", customerId, orderId)
}
