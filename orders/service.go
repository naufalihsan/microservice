package main

import (
	"context"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
)

type Service struct {
	store OrderStore
}

func NewService(store OrderStore) *Service {
	return &Service{store}
}

func (s *Service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	products, err := s.ValidateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	order := &pb.Order{
		Id:         "1",
		CustomerId: req.CustomerId,
		Products:   products,
		Status:     "pending",
	}

	return order, nil
}

func (s *Service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Product, error) {
	// temp inmem price id
	priceTable := map[string]string{
		"1": "price_1PM586FsDc9cxjmW18I8bBu3",
		"2": "price_1PM59AFsDc9cxjmW6TCaCDfZ",
	}

	if len(req.OrderProducts) == 0 {
		return nil, common.ErrEmptyOrderProducts
	}

	for _, orderProduct := range req.OrderProducts {
		if orderProduct.Quantity <= 0 {
			return nil, common.ErrInvalidOrderProductQuantity
		}
	}

	orderProducts := mergeOrderProductQuantity(req.OrderProducts)

	products := []*pb.Product{}
	for _, orderProduct := range orderProducts {
		products = append(products, &pb.Product{
			PriceId:  priceTable[orderProduct.ProductId],
			Quantity: orderProduct.Quantity,
		})
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
