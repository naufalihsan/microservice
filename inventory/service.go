package main

import (
	"context"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
)

type Service struct {
	store InventoryStore
}

func NewService(store InventoryStore) *Service {
	return &Service{store}
}

func (s *Service) Get(ctx context.Context, ids []string) ([]*pb.Product, error) {
	return s.store.Get(ctx, ids)
}

func (s *Service) Validate(ctx context.Context, orderProducts []*pb.OrderProduct) ([]*pb.Product, error) {
	productIds := make([]string, 0, len(orderProducts))
	for _, orderProduct := range orderProducts {
		productIds = append(productIds, orderProduct.ProductId)
	}

	inventories, err := s.store.Get(ctx, productIds)
	if err != nil {
		return nil, err
	}

	inventoryMap := make(map[string]*pb.Product, len(inventories))
	for _, inventory := range inventories {
		inventoryMap[inventory.Id] = inventory
	}

	products := make([]*pb.Product, 0, len(orderProducts))
	for _, orderProduct := range orderProducts {
		inventory, ok := inventoryMap[orderProduct.ProductId]
		if !ok || inventory.Quantity < orderProduct.Quantity {
			return nil, common.ErrProductOutOfStock
		}

		products = append(products, &pb.Product{
			Id:       inventory.Id,
			Name:     inventory.Name,
			PriceId:  inventory.PriceId,
			Quantity: orderProduct.Quantity,
		})
	}

	return products, nil
}

func (s *Service) Update(ctx context.Context, products []*pb.Product) error {
	productIds := make([]string, 0, len(products))
	for _, product := range products {
		productIds = append(productIds, product.Id)
	}

	inventories, err := s.store.Get(ctx, productIds)
	if err != nil {
		return err
	}

	inventoryMap := make(map[string]*pb.Product, len(inventories))
	for _, inventory := range inventories {
		inventoryMap[inventory.Id] = inventory
	}

	for _, product := range products {
		inventory, ok := inventoryMap[product.Id]
		if !ok {
			return common.ErrInvalidProduct
		}

		if inventory.Quantity < product.Quantity {
			return common.ErrProductOutOfStock
		}

		inventory.Quantity -= product.Quantity
	}

	return s.store.Update(ctx, inventoryMap)
}
