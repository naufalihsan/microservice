package main

import (
	"context"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
)

type Store struct {
	products map[string]*pb.Product
}

func NewStore() *Store {
	// mock mongodb (for simplicity)
	return &Store{
		products: map[string]*pb.Product{
			"1": {
				Id:       "1",
				Name:     "Fuji Apple",
				Quantity: 10,
				PriceId:  "price_1PM586FsDc9cxjmW18I8bBu3",
			},
			"2": {
				Id:       "2",
				Name:     "Honey Pineapple",
				Quantity: 20,
				PriceId:  "price_1PM59AFsDc9cxjmW6TCaCDfZ",
			},
		},
	}
}

func (s *Store) Get(ctx context.Context, ids []string) ([]*pb.Product, error) {
	products := make([]*pb.Product, 0, len(ids))
	for _, id := range ids {
		product, ok := s.products[id]
		if !ok {
			return nil, common.ErrInvalidProduct
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *Store) Update(ctx context.Context, inventoryMap map[string]*pb.Product) error {
	for _, inventory := range inventoryMap {
		s.products[inventory.Id] = inventory
	}
	return nil
}
