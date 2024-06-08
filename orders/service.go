package main

import (
	"context"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/naufalihsan/msvc-orders/gateway"
)

type Service struct {
	store   OrderStore
	gateway gateway.InventoryGateaway
}

func NewService(store OrderStore, gateway gateway.InventoryGateaway) *Service {
	return &Service{store, gateway}
}

func (s *Service) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, req.CustomerId, req.OrderId)
}

func (s *Service) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	return s.store.Update(ctx, order)
}

func (s *Service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	products, err := s.ValidateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	order, err := s.store.Create(ctx, req, products)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Product, error) {
	if len(req.OrderProducts) == 0 {
		return nil, common.ErrEmptyOrderProducts
	}

	for _, orderProduct := range req.OrderProducts {
		if orderProduct.Quantity <= 0 {
			return nil, common.ErrInvalidOrderProductQuantity
		}
	}

	orderProducts := mergeOrderProductQuantity(req.OrderProducts)

	products, err := s.gateway.Validate(ctx, orderProducts)
	if err != nil {
		return nil, err
	}

	return products, nil
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
