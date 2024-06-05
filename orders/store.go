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
		Status:     common.OrderStatusPending,
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

func (s *Store) Update(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	uniqueKey := GenerateUniqueKey(order.CustomerId, order.Id)

	updateOrder, ok := inmemStore[uniqueKey]
	if !ok {
		return nil, common.ErrInvalidOrder
	}

	if order.Status != "" {
		updateOrder.Status = order.Status
	}

	if order.PaymentLink != "" {
		updateOrder.PaymentLink = order.PaymentLink
	}

	inmemStore[uniqueKey] = updateOrder

	return updateOrder, nil
}

func GenerateUniqueKey(customerId, orderId string) string {
	return fmt.Sprintf("%s-%s", customerId, orderId)
}
