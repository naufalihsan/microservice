package main

import (
	"context"
	"log"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
)

type Service struct {
	store OrderStore
}

func NewService(store OrderStore) *Service {
	return &Service{store}
}

func (s *Service) CreateOrder(context.Context) error {
	return nil
}

func (s *Service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error {
	if len(req.OrderProducts) == 0 {
		return common.ErrEmptyOrderProducts
	}

	for _, orderProduct := range req.OrderProducts {
		if orderProduct.Quantity <= 0 {
			return common.ErrInvalidOrderProductQuantity
		}
	}

	orderProducts := mergeOrderProductQuantity(req.OrderProducts)
	log.Println(orderProducts)

	return nil
}

func mergeOrderProductQuantity(orderProducts []*pb.OrderProduct) []*pb.OrderProduct {
	// Create a map to keep track of the quantities for each ProductId
	productMap := make(map[string]*pb.OrderProduct)

	// Iterate over each order product
	for _, orderProduct := range orderProducts {
		if existingProduct, exists := productMap[orderProduct.ProductId]; exists {
			// If the product already exists in the map, update the quantity
			existingProduct.Quantity += orderProduct.Quantity
		} else {
			// If the product does not exist in the map, add it
			productMap[orderProduct.ProductId] = orderProduct
		}
	}

	// Convert the map back to a slice
	uniqueOrderProducts := make([]*pb.OrderProduct, 0, len(productMap))
	for _, orderProduct := range productMap {
		uniqueOrderProducts = append(uniqueOrderProducts, orderProduct)
	}

	return uniqueOrderProducts
}
